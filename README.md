[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![Work In Progress](https://img.shields.io/badge/Status-Work%20In%20Progress-yellow)](https://guide.unitvectorylabs.com/bestpractices/status/#work-in-progress)
 [![Go Report Card](https://goreportcard.com/badge/github.com/UnitVectorY-Labs/badgeindexer)](https://goreportcard.com/report/github.com/UnitVectorY-Labs/badgeindexer)

# badgeindexer

Collects README badges across an organization and generates a structured dashboard with per-repo and per-badge detail pages.

## Usage

### Crawl Command

Fetches README badges from GitHub repositories:

```bash
./badgeindexer -crawl -org <organization> [flags]
```

Flags:
- `-org <name>`: GitHub owner/organization name (required)
- `-private`: Include private repositories (default: public only)
- `-output <path>`: Directory for JSON output (default: `data`)

Requirements:
- `GITHUB_TOKEN` environment variable with a valid GitHub personal access token

Example:

```bash
export GITHUB_TOKEN=$(gh auth token)
./badgeindexer -crawl -org UnitVectorY-Labs
```

### Generate Command

Creates static HTML pages from crawl JSON data:

```bash
./badgeindexer -generate [flags]
```

Flags:
- `-output <path>`: Directory containing JSON data (default: `data`)
- `-html <path>`: Directory for HTML output (default: `output`)

Example:

```bash
./badgeindexer -generate
```

#### Development Mode

For development, you can override the embedded templates to load from disk instead. This allows live editing of templates and CSS without rebuilding the binary:

```bash
export TEMPLATE_PATH=./templates
./badgeindexer -generate
```

When `TEMPLATE_PATH` is set, templates and the `style.css` file are loaded from the specified directory instead of the embedded filesystem that is part of the binary.

## Configuration

### badges.json

Create a `badges.json` file to map badge patterns to human-readable names and categories:

```json
{
  "badges": [
    {
      "pattern": "https://img.shields.io/badge/license-MIT-blue.svg",
      "name": "MIT",
      "category": "License"
    },
    {
        "pattern": "https://img.shields.io/badge/Status-Work%20In%20Progress-yellow",
        "name": "Work In Progress",
        "category": "Project Status"
    }
  ]
}
```

Fields:
- `pattern`: URL pattern to match (supports `{ORG}`, `{REPO}`, and `.*` wildcards)
- `name`: Display name for the badge
- `category`: Category for grouping badges
- `placeholder`: Optional static image URL to display instead of the dynamic badge for a specific repository

Unrecognized badges are assigned to the "Unknown" category.

## Output

### JSON Output (from crawl)

Each repository creates a JSON file in the output directory:

```json
{
  "repository": "example-repo",
  "default_branch": "main",
  "readme_found": true,
  "badges": [
    {
      "alt_text": "License",
      "image_url": "https://img.shields.io/badge/License-MIT-blue.svg",
      "target_url": "https://opensource.org/licenses/MIT",
      "host_image": "img.shields.io",
      "host_target": "opensource.org"
    }
  ]
}
```

A `timestamp.json` file records the last crawl time.
