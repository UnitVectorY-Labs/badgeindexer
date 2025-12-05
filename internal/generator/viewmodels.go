package generator

import "github.com/UnitVectorY-Labs/badgeindexer/internal/models"

// DashboardViewModel is used for the index page.
type DashboardViewModel struct {
	OrgName          string
	TotalRepos       int
	TotalBadges      int
	ReposWithBadges  int
	ReposNoBadges    int
	Repositories     []RepoSummary
	BadgesByCategory []BadgeCategory
	UniqueBadgeCount int
	LastUpdated      string
}

// BadgeCategory groups badges by category.
type BadgeCategory struct {
	Name   string
	Badges []BadgeSummary
}

// RepoSummary is a summary of a repository for listing.
type RepoSummary struct {
	Name       string
	BadgeCount int
	BadgeIDs   []string // IDs of badges this repo has, for filtering
}

// BadgeSummary is a summary of a unique badge pattern.
type BadgeSummary struct {
	ImageURL string
	Pattern  string
	Name     string
	Category string
	Count    int
	ID       string
}

// RepoBadge represents a badge with its looked-up name and category.
type RepoBadge struct {
	ImageURL  string
	TargetURL string
	AltText   string
	Name      string
	Category  string
	ID        string
}

// RepoPageViewModel is used for individual repository pages.
type RepoPageViewModel struct {
	OrgName     string
	Repository  models.RepositoryData
	Badges      []RepoBadge
	LastUpdated string
}

// BadgeRepoBadge represents a badge instance for a specific repo on badge pages.
type BadgeRepoBadge struct {
	RepoName  string
	ImageURL  string
	TargetURL string
}

// BadgePageViewModel is used for individual badge pages.
type BadgePageViewModel struct {
	OrgName      string
	ImageURL     string
	Pattern      string
	Name         string
	Category     string
	Repositories []BadgeRepoBadge
	LastUpdated  string
}

// BadgeConfig represents the badges.json configuration.
type BadgeConfig struct {
	Badges []BadgeConfigEntry `json:"badges"`
}

// BadgeConfigEntry represents a single badge configuration entry.
type BadgeConfigEntry struct {
	ID          string `json:"id"`
	Pattern     string `json:"pattern"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Placeholder string `json:"placeholder,omitempty"`
}
