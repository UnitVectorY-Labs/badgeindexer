package crawler

import (
	"testing"
	"testing/fstest"
)

func TestLoadBadgeDomains(t *testing.T) {
	t.Parallel()

	fsys := fstest.MapFS{
		"badge-domains.yaml": {
			Data: []byte("domains:\n  - img.shields.io\n  - goreportcard.com\n"),
		},
	}

	domains, err := LoadBadgeDomains(fsys, "badge-domains.yaml")
	if err != nil {
		t.Fatalf("LoadBadgeDomains() error = %v", err)
	}

	if _, ok := domains["img.shields.io"]; !ok {
		t.Fatalf("expected img.shields.io to be loaded")
	}

	if _, ok := domains["goreportcard.com"]; !ok {
		t.Fatalf("expected goreportcard.com to be loaded")
	}
}
