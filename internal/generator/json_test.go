package generator

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRepositoryJSONRejectsInvalidInput(t *testing.T) {
	for name, input := range map[string]string{
		"duplicate names": `{"repository":"one","repository":"two"}`,
		"invalid UTF-8":   "{\"repository\":\"\xff\"}",
		"trailing value":  `{"repository":"one"} {}`,
	} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "repo.json"), []byte(input), 0600); err != nil {
				t.Fatal(err)
			}
			if err := Run(dir, t.TempDir(), embed.FS{}); err == nil || !strings.Contains(err.Error(), "failed to decode") {
				t.Fatalf("expected decode error, got %v", err)
			}
		})
	}
}

func TestBadgeConfigJSON(t *testing.T) {
	for _, tc := range []struct {
		name, input string
		wantName    string
	}{
		{"valid", `{"badges":[{"name":"Build"}]}`, "Build"},
		{"case sensitive", `{"Badges":[{"name":"Build"}]}`, ""},
		{"duplicate", `{"badges":[],"badges":[{"name":"Build"}]}`, ""},
		{"trailing", `{"badges":[{"name":"Build"}]} {}`, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "badges.json")
			if err := os.WriteFile(path, []byte(tc.input), 0600); err != nil {
				t.Fatal(err)
			}
			config := loadBadgeConfig(path)
			got := ""
			if len(config.Badges) > 0 {
				got = config.Badges[0].Name
			}
			if got != tc.wantName {
				t.Fatalf("got name %q, want %q", got, tc.wantName)
			}
		})
	}
}
