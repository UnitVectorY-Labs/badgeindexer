# UVY Design adoption

The interface follows the current [UVY Design guidance](https://github.com/UnitVectorY-Labs/uvy-design) (reviewed at commit `64017af`). The guidance is a human-readable reference, not a versioned dependency. BadgeIndexer owns its implementation and serves interface assets locally.

The project title and organization subtitle identify the application and its data scope. The theme toggle follows the operating system until a choice is saved in the `uvy-theme` cookie. The dashboard, filtering, repository and badge pages, HTMX navigation, and footer retain their established behavior.

The theme controller and color variables originated in the former UVY Design v0.1.0 asset package and are maintained locally under its MIT license (`assets/UVY-DESIGN-LICENSE`). The inline moon and sun are Tabler icons (MIT, `assets/LICENSE`). HTMX 4.0.0 is served locally under its BSD license (`assets/HTMX-LICENSE`), with the same SHA-384 integrity value as the previously hosted copy.

Badge artwork is indexed crawl data and remains loaded from its original external URLs. Dashboard filter state resets when navigating back to the dashboard, as before.
