package crawler

import "testing"

func TestExtractBadges(t *testing.T) {
	t.Parallel()

	detector := badgeDetector{
		domains: map[string]struct{}{
			"archestra.ai":     {},
			"codecov.io":       {},
			"goreportcard.com": {},
			"img.shields.io":   {},
			"javadoc.io":       {},
		},
	}

	tests := []struct {
		name    string
		content string
		want    int
	}{
		{
			name:    "extracts shields badge",
			content: `[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)`,
			want:    1,
		},
		{
			name:    "extracts shields badge without badge path segment",
			content: `[![Maven Central](https://img.shields.io/maven-central/v/com.unitvectory/consistgen)](https://central.sonatype.com/artifact/com.unitvectory/consistgen)`,
			want:    1,
		},
		{
			name: "ignores linked screenshot",
			content: `# gowebshot

[![gowebshot CLI](./docs/cli.png)](cli.png)
`,
			want: 0,
		},
		{
			name:    "extracts html badge",
			content: `<a href="https://github.com/example/project/actions/workflows/test.yml"><img src="https://github.com/example/project/actions/workflows/test.yml/badge.svg" alt="Tests"></a>`,
			want:    1,
		},
		{
			name:    "extracts local badge when filename says badge",
			content: `[![Build status](./docs/build-badge.svg)](https://example.com/build)`,
			want:    1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			badges := extractBadges([]byte(tt.content), detector)
			if len(badges) != tt.want {
				t.Fatalf("extractBadges() returned %d badges, want %d", len(badges), tt.want)
			}
		})
	}
}
