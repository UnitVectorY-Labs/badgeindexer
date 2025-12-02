package crawler

import (
	"net/url"
	"regexp"

	"github.com/UnitVectorY-Labs/badgeindexer/internal/models"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// extractBadges parses the README content and returns a list of badges.
func extractBadges(content []byte) []models.Badge {
	var badges []models.Badge

	// 1. Parse Markdown AST
	md := goldmark.New()
	reader := text.NewReader(content)
	doc := md.Parser().Parse(reader)

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		// Check for Link nodes
		if link, ok := n.(*ast.Link); ok {
			// Check if the link contains an Image
			for child := link.FirstChild(); child != nil; child = child.NextSibling() {
				if img, ok := child.(*ast.Image); ok {
					badge := models.Badge{
						AltText:   string(img.Text(content)),
						ImageURL:  string(img.Destination),
						TargetURL: string(link.Destination),
					}
					normalizeBadge(&badge)
					badges = append(badges, badge)
				}
			}
		}
		return ast.WalkContinue, nil
	})

	// 2. Regex fallback for HTML badges: <a href="..."><img src="..." alt="..."></a>
	// This is a simple regex and might not catch all edge cases, but covers the standard pattern.
	htmlBadgeRegex := regexp.MustCompile(`<a\s+href="([^"]+)"[^>]*>\s*<img\s+src="([^"]+)"(?:\s+alt="([^"]*)")?[^>]*>\s*</a>`)
	matches := htmlBadgeRegex.FindAllSubmatch(content, -1)
	for _, match := range matches {
		badge := models.Badge{
			TargetURL: string(match[1]),
			ImageURL:  string(match[2]),
			AltText:   string(match[3]),
		}
		normalizeBadge(&badge)

		badges = append(badges, badge)
	}

	return badges
}

func normalizeBadge(b *models.Badge) {
	if u, err := url.Parse(b.ImageURL); err == nil {
		b.HostImage = u.Host
	}
	if u, err := url.Parse(b.TargetURL); err == nil {
		b.HostTarget = u.Host
	}
}
