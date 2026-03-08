package crawler

import (
	"fmt"
	"io/fs"
	"strings"

	"gopkg.in/yaml.v3"
)

type badgeDomainConfig struct {
	Domains []string `yaml:"domains"`
}

func LoadBadgeDomains(fsys fs.FS, path string) (map[string]struct{}, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("read badge domains: %w", err)
	}

	var config badgeDomainConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse badge domains: %w", err)
	}

	domains := make(map[string]struct{}, len(config.Domains))
	for _, domain := range config.Domains {
		domain = strings.ToLower(strings.TrimSpace(domain))
		if domain == "" {
			continue
		}
		domains[domain] = struct{}{}
	}

	if len(domains) == 0 {
		return nil, fmt.Errorf("badge domains config %q is empty", path)
	}

	return domains, nil
}
