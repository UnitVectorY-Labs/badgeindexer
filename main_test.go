package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/UnitVectorY-Labs/badgeindexer/internal/generator"
)

func TestGeneratedSiteAssetsAndFragments(t *testing.T) {
	input, output := t.TempDir(), t.TempDir()
	fixture := `{"repository":"example","repository_url":"https://github.com/Test/example","default_branch":"main","readme_found":true,"badges":[{"image_url":"https://img.shields.io/badge/test-ok-green","target_url":"https://example.com","alt_text":"test"}]}`
	if err := os.WriteFile(filepath.Join(input, "example.json"), []byte(fixture), 0600); err != nil {
		t.Fatal(err)
	}
	if err := generator.Run(input, output, appFS); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"style.css", "uvy-tokens.css", "uvy-theme.js", "htmx-4.0.0.min.js"} {
		embedded, err := appFS.ReadFile("templates/" + name)
		if err != nil {
			t.Fatal(err)
		}
		generated, err := os.ReadFile(filepath.Join(output, name))
		if err != nil {
			t.Fatal(err)
		}
		if string(embedded) != string(generated) {
			t.Errorf("asset %s differs", name)
		}
	}
	badges, err := filepath.Glob(filepath.Join(output, "badges", "*.html"))
	if err != nil || len(badges) != 1 {
		t.Fatalf("badge pages: %v %v", badges, err)
	}
	for _, path := range []string{filepath.Join(output, "index.html"), filepath.Join(output, "repos", "example.html"), badges[0]} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		s := string(data)
		for _, forbidden := range []string{"unpkg.com", "cdnjs.cloudflare.com"} {
			if strings.Contains(s, forbidden) {
				t.Errorf("%s contains %s", path, forbidden)
			}
		}
		theme, tokens, style, toggle, main := strings.Index(s, "/uvy-theme.js"), strings.Index(s, "/uvy-tokens.css"), strings.Index(s, "/style.css"), strings.Index(s, "data-uvy-theme-toggle"), strings.Index(s, `<main id="content">`)
		if theme < 0 || !(theme < tokens && tokens < style && toggle < main && toggle > 0) {
			t.Errorf("asset order or toggle placement in %s", path)
		}
		rel, _ := filepath.Rel(output, path)
		snippet, err := os.ReadFile(filepath.Join(output, "snippets", rel))
		if err != nil {
			t.Fatal(err)
		}
		fragment := strings.TrimSpace(string(snippet))
		if strings.Contains(fragment, "<head>") || !strings.Contains(s, fragment) {
			t.Errorf("fragment mismatch: %s", rel)
		}
	}
}
