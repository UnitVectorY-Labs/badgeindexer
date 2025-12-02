The badgeindexer project is a Go command line application that has two main functions:
1. Crawling a specified GitHub organization to find repositories containing badges in their README.md file and storing the results in JSON files.
2. Generating static HTML pages from the collected badge data using Go templates for viewing in a web browser.

For development usage the `data` and `output` directories are used to store the crawled data and generated HTML files respectively. These directories are created and populated when the application is run and are not part of the application itself.

The framework uses HTMX for dynamic navigation with a single-page application approach. A single `base.html` template serves as the foundation, with `index.html` as the main entry point and snippet templates providing content that is dynamically loaded via HTMX. All styles are defined in `assets/style.css`. The HTML templates are located in the `templates` directory:

- `base.html`: Base template with full HTML structure (header, footer, scripts)
- `index.html`: Main entry page that uses base template with dashboard content
- `repo.html`: Repository detail page that uses base template with repo snippet
- `badge.html`: Badge detail page that uses base template with badge snippet
- `index_snippet.html`: Dashboard content snippet (stats, badge filters, repository list)
- `badge_snippet.html`: Individual badge detail snippet
- `repo_snippet.html`: Individual repository detail snippet

The generated output includes:
- `output/index.html`: The main entry point (full HTML page)
- `output/repos/{name}.html`: Full HTML pages for each repository
- `output/badges/{id}.html`: Full HTML pages for each badge
- `output/snippets/`: Content snippets loaded by HTMX for dynamic navigation