package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/UnitVectorY-Labs/badgeindexer/internal/crawler"
	"github.com/UnitVectorY-Labs/badgeindexer/internal/generator"
)

// templateFS embeds all HTML templates from the templates directory.
//
//go:embed templates/*.html
//go:embed templates/style.css
var templateFS embed.FS

func main() {
	crawlMode := flag.Bool("crawl", false, "Run the crawler phase")
	genMode := flag.Bool("generate", false, "Run the generator phase")
	orgName := flag.String("org", "", "GitHub Organization name (required for crawl)")
	includePrivate := flag.Bool("private", false, "Include private repositories (default: public only)")
	outputDir := flag.String("output", "data", "Directory for data output (crawl) or input (generate)")
	htmlDir := flag.String("html", "output", "Directory for HTML output (generate)")

	flag.Parse()

	if *crawlMode && *genMode {
		fmt.Println("Error: Cannot run both -crawl and -generate at the same time.")
		os.Exit(1)
	}

	if !*crawlMode && !*genMode {
		fmt.Println("Usage: badge-indexer [-crawl | -generate] [options]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *crawlMode {
		if *orgName == "" {
			fmt.Println("Error: -org is required for crawl mode.")
			os.Exit(1)
		}
		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			fmt.Println("Error: GITHUB_TOKEN environment variable is required for crawl mode.")
			os.Exit(1)
		}
		if err := crawler.Run(*orgName, *outputDir, token, *includePrivate); err != nil {
			fmt.Printf("Crawl failed: %v\n", err)
			os.Exit(1)
		}
	}

	if *genMode {
		if err := generator.Run(*outputDir, *htmlDir, templateFS); err != nil {
			fmt.Printf("Generation failed: %v\n", err)
			os.Exit(1)
		}
	}
}
