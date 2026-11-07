# UVY Design adoption

Pinned to [UVY Design v0.1.0](https://github.com/UnitVectorY-Labs/uvy-design/tree/12f8e3d964785d78628346516351fe4dac583cef), source commit `12f8e3d964785d78628346516351fe4dac583cef`. Nothing is fetched from the design repository at runtime.

## Copied files

Shared files are copied verbatim into `templates/`, embedded by `main.go` and written to the generated site root by `internal/generator`.

| Path | Source in v0.1.0 | SHA-256 | License |
| --- | --- | --- | --- |
| `templates/uvy-tokens.css` | `assets/uvy-tokens.css` | `0e97baaece886e0015c4238eaf7dbed9ea68f330c0f50729311c43b4a5a98533` | MIT, `assets/UVY-DESIGN-LICENSE` |
| `templates/uvy-theme.js` | `assets/uvy-theme.js` | `b105c0ddec8c19b6a380a33ffa5f9f558ca49baa477921a1dee9596d77bc61ed` | MIT, `assets/UVY-DESIGN-LICENSE` |
| `templates/htmx-4.0.0.min.js` | `vendor/htmx-4.0.0.min.js` | `e484d9171a9db30a39c8f16e3d709d4137f3211c659f8e6125816635033d593f` | BSD, `assets/HTMX-LICENSE` |

`uvy-elements.css` is not used. The HTMX file is the same version (4.0.0) the site previously loaded from unpkg; its SHA-384 is `BvJpBiO8Kh31EqtJe5DRIeWrHWnCGkwytKs9NKFi86Hhw96dEqdEMzZDeK9iEGTc`, identical to the previous `integrity` value, which is kept on the local script tag.

The inline moon and sun in the theme toggle are the release's Tabler outline icons (`icons/moon.svg` SHA-256 `7cabaae40eea65310da478d3cad30b23cf1d7e99adad56378f82824338079668`, `icons/sun.svg` SHA-256 `6677490605af8987d55d6e4185b1f6c241eced65b63ccfece68cadb7ce90fbb5`), MIT licensed in `assets/LICENSE`.

## Adapter

- Head order in `base.html`, `index.html`, `repo.html` and `badge.html`: synchronous `uvy-theme.js`, `uvy-tokens.css`, `style.css`, then HTMX.
- The theme toggle (`data-uvy-theme-toggle`, initially `hidden`) sits at the top right of `<header>`, outside the `#content` swap region.
- `templates/style.css` maps the existing header, stat cards, badge filters, search, repository/badge tables, info boxes and footer onto the shared tokens at standard density. Layout, navigation, filter behavior and content order are unchanged. The old uppercase category labels are sentence case per the style guide.
- Minimal accessibility repair: the repository search input has `aria-label="Search repositories"`.
- HTMX 4 is unchanged; `hx-get`, `hx-target` and `hx-push-url` navigation is untouched.

## Approved exceptions and notes

- Badge artwork (shields.io, codecov.io, javadoc.io, etc.) is indexed crawl data rendered by URL, not an application asset. It remains remote; a build-time snapshot would be a separate change.
- The footer retains its existing organization, MIT and repository links.
- Dashboard filter state still resets when navigating back to the dashboard, as before; persisting it would be a new behavior.

## Verification and updates

Reset an explicit theme by clearing the site's `uvy-theme` cookie and reloading. To update, review the new release diff and licenses, copy the shared files verbatim, update the hashes above, then run `go test ./...`, `go vet ./...` and `go build ./...`, generate a site from real crawl data and compare browser captures and interactions in both themes at desktop, 390px and 320px widths.
