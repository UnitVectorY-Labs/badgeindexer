package models

import (
	"encoding/json/v2"
	"strings"
	"testing"
)

func TestRepositoryDataJSONEmptyBadges(t *testing.T) {
	data, err := json.Marshal(RepositoryData{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"badges":[]`) {
		t.Fatalf("nil badges should encode as an empty array: %s", data)
	}
	for _, input := range []string{`{"badges":null}`, `{"badges":[]}`} {
		var repo RepositoryData
		if err := json.Unmarshal([]byte(input), &repo); err != nil || len(repo.Badges) != 0 {
			t.Fatalf("reading %s: badges=%v, err=%v", input, repo.Badges, err)
		}
	}
}
