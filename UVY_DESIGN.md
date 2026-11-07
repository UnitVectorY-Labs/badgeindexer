# UVY Design adoption

Pinned to [UVY Design v0.1.0](https://github.com/UnitVectorY-Labs/uvy-design/releases/tag/v0.1.0), source commit `396a65021c2bcab8dbcb3d1bb4cb498c807a4396`. The design assets are MIT licensed; see `assets/UVY-DESIGN-LICENSE`. SHA-256 of verbatim copies:

- `templates/uvy-tokens.css`: `0e97baaece886e0015c4238eaf7dbed9ea68f330c0f50729311c43b4a5a98533`
- `templates/uvy-theme.js`: `b105c0ddec8c19b6a380a33ffa5f9f558ca49baa477921a1dee9596d77bc61ed`
- `templates/htmx-4.0.0.min.js`: `e484d9171a9db30a39c8f16e3d709d4137f3211c659f8e6125816635033d593f`, HTMX BSD license in `assets/HTMX-LICENSE`.

The inline moon and sun are Tabler icons distributed in the release, MIT licensed in `assets/LICENSE`. Source SVG SHA-256: moon `7cabaae40eea65310da478d3cad30b23cf1d7e99adad56378f82824338079668`, sun `6677490605af8987d55d6e4185b1f6c241eced65b63ccfece68cadb7ce90fbb5`.

`templates/style.css` maps the existing header, sections, filters, tables and statistics to the shared tokens at standard density. The theme toggle remains outside HTMX swaps. A session-storage adapter retains dashboard filters through navigation. Clear the `uvy-theme` cookie and reload to restore OS preference. Badge images are indexed external content, not locally vendored application assets. Existing navigation and footer are retained.

To update, review the new design release and license changes, copy the shared assets verbatim, update the hashes here, then run `go test ./...`, `go vet ./...`, `go build ./...` and compare browser captures and interactions in both themes and at mobile widths. Do not fetch design assets at runtime.
