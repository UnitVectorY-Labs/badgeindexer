package models

// Badge represents a single badge found in a README.
type Badge struct {
	AltText    string `json:"alt_text"`
	ImageURL   string `json:"image_url"`
	TargetURL  string `json:"target_url"`
	HostImage  string `json:"host_image"`
	HostTarget string `json:"host_target"`
}

// RepositoryData represents the crawled data for a single repository.
type RepositoryData struct {
	Repository    string  `json:"repository"`
	RepositoryURL string  `json:"repository_url"`
	DefaultBranch string  `json:"default_branch"`
	ReadmeFound   bool    `json:"readme_found"`
	Badges        []Badge `json:"badges"`
}
