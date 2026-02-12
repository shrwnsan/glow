# Changelog

All notable changes to Glow will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Superscript and subscript rendering support for HTML `<sup>` and `<sub>` tags
  - Converts HTML tags to Unicode superscript/subscript characters before rendering
  - Supports mathematical expressions like `E = mc²` (E = mc<sup>2</sup>)
  - Supports chemical formulas like `H₂O` (H<sub>2</sub>O)
  - Renders powers like `x² + y² = r²` (x<sup>2</sup> + y<sup>2</sup> = r<sup>2</sup>)

### Changed
- Bump Go installation path to v2 in README

### Fixed
- No fixes in this release

## [2.1.1] - 2025-12-15

### Added
- Initial release

[2.1.1]: https://github.com/charmbracelet/glow/releases/tag/v2.1.1
