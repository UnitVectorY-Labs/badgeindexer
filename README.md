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

Note: Archived repositories are always excluded from crawling, as they cannot be modified and are treated as if they do not exist.

Requirements:
- `GITHUB_TOKEN` environment variable with a valid GitHub personal access token

Badge detection behavior:
- Linked images are treated as badges when the image URL, target URL, or image filename contains `badge`
- Linked images are also treated as badges when the image host appears in `badge-domains.yaml`
- The badge domain list is embedded into the binary through Go's `embed.FS`, so crawl behavior is consistent at runtime without requiring an external file next to the executable

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

### badge-domains.yaml

`badge-domains.yaml` lives at the repository root and defines image hostnames that should always be considered badge providers during crawl. This file is embedded into the binary and loaded from the virtual filesystem at startup.

```yaml
domains:
  - img.shields.io
  - goreportcard.com
  - codecov.io
```

If a linked image does not come from one of these configured hosts, it must still include `badge` in the image URL, target URL, or filename to be treated as a badge. This prevents normal linked screenshots and other README images from being misclassified.

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
