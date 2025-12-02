package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/UnitVectorY-Labs/badgeindexer/internal/models"
	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

const DefaultWorkerCount = 10

// Run executes the crawl phase.
func Run(orgName, outputDir, token string, includePrivate bool) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 1. List all repositories
	fmt.Printf("Fetching repositories for org: %s...\n", orgName)
	var allRepos []*github.Repository
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, orgName, opt)
		if err != nil {
			return fmt.Errorf("failed to list repositories: %w", err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	fmt.Printf("Found %d repositories.\n", len(allRepos))

	// Filter out private repos unless includePrivate is set
	if !includePrivate {
		var publicRepos []*github.Repository
		for _, repo := range allRepos {
			if !repo.GetPrivate() {
				publicRepos = append(publicRepos, repo)
			}
		}
		fmt.Printf("Filtered to %d public repositories.\n", len(publicRepos))
		allRepos = publicRepos
	}

	// 2. Worker Pool for fetching READMEs
	jobs := make(chan *github.Repository, len(allRepos))
	results := make(chan error, len(allRepos))
	var wg sync.WaitGroup

	// Concurrency limit
	workerCount := DefaultWorkerCount
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for repo := range jobs {
				results <- processRepo(ctx, client, repo, outputDir)
			}
		}()
	}

	for _, repo := range allRepos {
		jobs <- repo
	}
	close(jobs)

	wg.Wait()
	close(results)

	// Check for errors
	errCount := 0
	for err := range results {
		if err != nil {
			fmt.Printf("Error processing repo: %v\n", err)
			errCount++
		}
	}

	fmt.Printf("Crawl complete. Errors: %d\n", errCount)

	// Write timestamp.json
	timestampData := map[string]string{
		"last_crawled": time.Now().Format(time.RFC3339Nano),
	}
	tsFile, err := os.Create(filepath.Join(outputDir, "timestamp.json"))
	if err != nil {
		return fmt.Errorf("failed to create timestamp.json: %w", err)
	}
	defer tsFile.Close()

	tsEncoder := json.NewEncoder(tsFile)
	tsEncoder.SetIndent("", "  ")
	if err := tsEncoder.Encode(timestampData); err != nil {
		return fmt.Errorf("failed to encode timestamp.json: %w", err)
	}

	return nil
}

func processRepo(ctx context.Context, client *github.Client, repo *github.Repository, outputDir string) error {
	repoName := repo.GetName()
	defaultBranch := repo.GetDefaultBranch()

	// Try to fetch README
	readme, _, err := client.Repositories.GetReadme(ctx, repo.GetOwner().GetLogin(), repoName, nil)

	readmeFound := false
	var badges []models.Badge

	if err == nil && readme != nil {
		readmeFound = true
		var contentStr string
		contentStr, err = readme.GetContent()
		if err != nil {
			return fmt.Errorf("failed to decode readme for %s: %w", repoName, err)
		}
		badges = extractBadges([]byte(contentStr))
	}

	data := models.RepositoryData{
		Repository:    repoName,
		RepositoryURL: repo.GetHTMLURL(),
		DefaultBranch: defaultBranch,
		ReadmeFound:   readmeFound,
		Badges:        badges,
	}

	// Save to JSON
	filename := filepath.Join(outputDir, fmt.Sprintf("%s.json", normalizeRepoName(repoName)))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func normalizeRepoName(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), "/", "-")
}
