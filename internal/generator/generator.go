package generator

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/UnitVectorY-Labs/badgeindexer/internal/models"
)

// Run executes the generation phase.
func Run(inputDir, outputDir string) error {
	fmt.Printf("Starting generation from: %s, output to: %s\n", inputDir, outputDir)

	// Ensure output directories exist
	for _, dir := range []string{
		filepath.Join(outputDir, "repos"),
		filepath.Join(outputDir, "badges"),
		filepath.Join(outputDir, "snippets"),
		filepath.Join(outputDir, "snippets", "repos"),
		filepath.Join(outputDir, "snippets", "badges"),
	} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Load badge configuration
	badgeConfig := loadBadgeConfig("badges.json")

	// Load timestamp
	lastUpdated := loadTimestamp(inputDir)

	// Load Data
	files, err := filepath.Glob(filepath.Join(inputDir, "*.json"))
	if err != nil {
		return fmt.Errorf("failed to list input files: %w", err)
	}

	var repos []models.RepositoryData
	var orgName string
	for _, f := range files {
		if strings.HasSuffix(f, "timestamp.json") {
			continue
		}

		file, err := os.Open(f)
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", f, err)
		}

		var repo models.RepositoryData
		if err := json.NewDecoder(file).Decode(&repo); err != nil {
			file.Close()
			return fmt.Errorf("failed to decode %s: %w", f, err)
		}
		file.Close()
		repos = append(repos, repo)

		// Extract org name from repository URL (e.g., https://github.com/OrgName/RepoName)
		if orgName == "" && repo.RepositoryURL != "" {
			if u, err := url.Parse(repo.RepositoryURL); err == nil {
				parts := strings.Split(strings.Trim(u.Path, "/"), "/")
				if len(parts) >= 1 {
					orgName = parts[0]
				}
			}
		}
	}

	// Build ViewModels
	dashboardVM, badgeMap := buildDashboard(repos, badgeConfig, orgName, lastUpdated)

	// Parse Templates
	funcMap := template.FuncMap{
		"urlize": normalizeRepoName,
	}
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	// Render Index (full page and snippet)
	if err := renderTemplate(tmpl, filepath.Join(outputDir, "index.html"), "index.html", dashboardVM); err != nil {
		return err
	}
	if err := renderTemplate(tmpl, filepath.Join(outputDir, "snippets", "index.html"), "index_snippet.html", dashboardVM); err != nil {
		return err
	}
	if err := renderTemplate(tmpl, filepath.Join(outputDir, "snippets", "home.html"), "index_snippet.html", dashboardVM); err != nil {
		return err
	}

	// Render Repo Pages
	for _, repo := range repos {
		// Build enhanced badge list with names/categories
		var repoBadges []RepoBadge
		for _, b := range repo.Badges {
			pattern := canonicalizeURL(b.ImageURL, orgName, repo.Repository)
			name, category, _, id := lookupBadge(pattern, badgeConfig)
			repoBadges = append(repoBadges, RepoBadge{
				ImageURL:  b.ImageURL,
				TargetURL: b.TargetURL,
				AltText:   b.AltText,
				Name:      name,
				Category:  category,
				ID:        id,
			})
		}

		vm := RepoPageViewModel{
			OrgName:     orgName,
			Repository:  repo,
			Badges:      repoBadges,
			LastUpdated: lastUpdated,
		}
		baseName := normalizeRepoName(repo.Repository)

		// Render full page
		if err := renderPageWithSnippet(tmpl, filepath.Join(outputDir, "repos", baseName+".html"), "repo.html", "repo_snippet.html", vm); err != nil {
			return err
		}
		// Render snippet
		if err := renderTemplate(tmpl, filepath.Join(outputDir, "snippets", "repos", baseName+".html"), "repo_snippet.html", vm); err != nil {
			return err
		}
	}

	// Render Badge Pages
	for pattern, info := range badgeMap {
		id := info.ID

		// Build repo list with badge info for each repo
		var repoBadges []BadgeRepoBadge
		for _, repoName := range info.Repos {
			// Find the original badge URL for this repo
			imageURL := info.SampleImage
			targetURL := ""
			for _, repo := range repos {
				if repo.Repository == repoName {
					for _, b := range repo.Badges {
						p := canonicalizeURL(b.ImageURL, orgName, repo.Repository)
						if p == pattern {
							imageURL = b.ImageURL
							targetURL = b.TargetURL
							break
						}
					}
					break
				}
			}
			repoBadges = append(repoBadges, BadgeRepoBadge{
				RepoName:  repoName,
				ImageURL:  imageURL,
				TargetURL: targetURL,
			})
		}

		vm := BadgePageViewModel{
			OrgName:      orgName,
			ImageURL:     info.SampleImage,
			Pattern:      pattern,
			Name:         info.Name,
			Category:     info.Category,
			Repositories: repoBadges,
			LastUpdated:  lastUpdated,
		}

		// Render full page
		if err := renderPageWithSnippet(tmpl, filepath.Join(outputDir, "badges", id+".html"), "badge.html", "badge_snippet.html", vm); err != nil {
			return err
		}
		// Render snippet
		if err := renderTemplate(tmpl, filepath.Join(outputDir, "snippets", "badges", id+".html"), "badge_snippet.html", vm); err != nil {
			return err
		}
	}

	// Copy Assets
	if err := copyFile("assets/style.css", filepath.Join(outputDir, "style.css")); err != nil {
		return fmt.Errorf("failed to copy style.css: %w", err)
	}

	fmt.Println("Generation complete.")
	return nil
}

type badgeInfo struct {
	SampleImage string
	Placeholder string
	Name        string
	Category    string
	ID          string
	Repos       []string
}

func buildDashboard(repos []models.RepositoryData, config *BadgeConfig, orgName, lastUpdated string) (DashboardViewModel, map[string]*badgeInfo) {
	vm := DashboardViewModel{
		OrgName:     orgName,
		TotalRepos:  len(repos),
		LastUpdated: lastUpdated,
	}

	badgeMap := make(map[string]*badgeInfo) // CanonicalPattern -> info

	var repoSummaries []RepoSummary
	for _, r := range repos {
		vm.TotalBadges += len(r.Badges)

		summary := RepoSummary{
			Name:       r.Repository,
			BadgeCount: len(r.Badges),
			BadgeIDs:   []string{},
		}

		if len(r.Badges) > 0 {
			vm.ReposWithBadges++
		} else {
			vm.ReposNoBadges++
		}

		for _, b := range r.Badges {
			pattern := canonicalizeURL(b.ImageURL, orgName, r.Repository)
			_, _, _, id := lookupBadge(pattern, config)
			summary.BadgeIDs = append(summary.BadgeIDs, id)

			if _, exists := badgeMap[pattern]; !exists {
				name, category, placeholder, id := lookupBadge(pattern, config)
				displayImage := b.ImageURL
				if placeholder != "" {
					displayImage = placeholder
				}
				badgeMap[pattern] = &badgeInfo{
					SampleImage: displayImage,
					Placeholder: placeholder,
					Name:        name,
					Category:    category,
					ID:          id,
					Repos:       []string{},
				}
			}
			badgeMap[pattern].Repos = append(badgeMap[pattern].Repos, r.Repository)
		}

		repoSummaries = append(repoSummaries, summary)
	}

	// Sort Repos
	sort.Slice(repoSummaries, func(i, j int) bool {
		return repoSummaries[i].Name < repoSummaries[j].Name
	})
	vm.Repositories = repoSummaries

	// Build Badges by Category
	categoryMap := make(map[string][]BadgeSummary)
	for pattern, info := range badgeMap {
		bs := BadgeSummary{
			ImageURL: info.SampleImage,
			Pattern:  pattern,
			Name:     info.Name,
			Category: info.Category,
			Count:    len(info.Repos),
			ID:       info.ID,
		}
		categoryMap[info.Category] = append(categoryMap[info.Category], bs)
	}

	// Sort categories and badges within
	var categories []string
	for cat := range categoryMap {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	// Move "Unknown" to end if present
	var hasUnknown bool
	var filtered []string
	for _, cat := range categories {
		if cat == "Unknown" {
			hasUnknown = true
		} else {
			filtered = append(filtered, cat)
		}
	}
	if hasUnknown {
		categories = append(filtered, "Unknown")
	} else {
		categories = filtered
	}

	for _, cat := range categories {
		badges := categoryMap[cat]
		sort.Slice(badges, func(i, j int) bool {
			return badges[i].Count > badges[j].Count
		})
		vm.BadgesByCategory = append(vm.BadgesByCategory, BadgeCategory{
			Name:   cat,
			Badges: badges,
		})
	}

	vm.UniqueBadgeCount = len(badgeMap)

	return vm, badgeMap
}

func loadBadgeConfig(path string) *BadgeConfig {
	file, err := os.Open(path)
	if err != nil {
		return &BadgeConfig{}
	}
	defer file.Close()

	var config BadgeConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return &BadgeConfig{}
	}
	return &config
}

func lookupBadge(pattern string, config *BadgeConfig) (name, category, placeholder, id string) {
	for _, entry := range config.Badges {
		// Convert the config pattern to a regex
		regexPattern := regexp.QuoteMeta(entry.Pattern)
		regexPattern = strings.ReplaceAll(regexPattern, `\{ORG\}`, `[^/]+`)
		regexPattern = strings.ReplaceAll(regexPattern, `\{REPO\}`, `[^/]+`)
		regexPattern = strings.ReplaceAll(regexPattern, `\.*`, `.*`)

		re, err := regexp.Compile("^" + regexPattern + "$")
		if err != nil {
			continue
		}

		if re.MatchString(pattern) {
			return entry.Name, entry.Category, entry.Placeholder, entry.ID
		}
	}
	return "Unknown", "Unknown", "", normalizeRepoName(pattern)
}

func canonicalizeURL(rawURL, org, repo string) string {
	result := rawURL
	lowerURL := strings.ToLower(rawURL)
	lowerOrg := strings.ToLower(org)
	lowerRepo := strings.ToLower(repo)

	// Replace Org
	if idx := strings.Index(lowerURL, lowerOrg); idx != -1 {
		result = result[:idx] + "{ORG}" + result[idx+len(org):]
		lowerURL = strings.ToLower(result)
	}

	// Replace Repo - always add /.* after {REPO} for consistency
	if idx := strings.Index(lowerURL, lowerRepo); idx != -1 {
		result = result[:idx] + "{REPO}/.*"
	}

	// Remove token query param if present
	u, err := url.Parse(result)
	if err == nil {
		q := u.Query()
		if q.Get("token") != "" {
			q.Del("token")
			if len(q) == 0 {
				u.RawQuery = ""
			} else {
				u.RawQuery = q.Encode()
			}
			result = u.String()
		}
	}

	return result
}

func renderTemplate(tmpl *template.Template, path, name string, data interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", path, err)
	}
	defer file.Close()
	return tmpl.ExecuteTemplate(file, name, data)
}

// renderPageWithSnippet renders a full page using a template that extends base.html
func renderPageWithSnippet(tmpl *template.Template, path, templateName, snippetName string, data interface{}) error {
	return renderTemplate(tmpl, path, templateName, data)
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func normalizeRepoName(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), "/", "-")
}

func hashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))[:12]
}

func loadTimestamp(inputDir string) string {
	tsPath := filepath.Join(inputDir, "timestamp.json")
	file, err := os.Open(tsPath)
	if err != nil {
		return "Unknown"
	}
	defer file.Close()

	var data map[string]string
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return "Unknown"
	}

	lastCrawled, ok := data["last_crawled"]
	if !ok {
		return "Unknown"
	}

	// Parse RFC3339Nano timestamp
	t, err := time.Parse(time.RFC3339Nano, lastCrawled)
	if err != nil {
		return "Unknown"
	}

	// Format as "January 2, 2006 15:04 MST"
	return t.UTC().Format("January 2, 2006 15:04 MST")
}
