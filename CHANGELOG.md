# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.1] - 2026-10-22

### Changed

- Require Go 1.25 or later to build the project, and use Go 1.27.3 for CI and release builds.
- Update the Go dependencies, including Resty, Cobra, and pflag.

### Fixed

- Update the release workflow to use GoReleaser v2, build snapshots for pull requests, and run on numeric version tags.

### Security

- Update `golang.org/x/net` to v0.55.0 to address a vulnerable dependency.

## [0.3.0] - 2022-08-14

### Added

- Add `search tweets` subcommands.

### Changed

- Rename the project and executable from `get-old-tweets` to `squawks`.
- Rename the `twitter` package to `api`.

## [0.2.1] - 2022-08-12

### Added

- Add `--near` and `--within` location search options.
- Add `--include` and `--exclude` tweet type options and allow multiple filters.

## [0.2.0] - 2022-08-10

### Added

- Add `--url` and `--geocode` search options.

### Changed

- Export location as separate `latitude` and `longitude` CSV columns.

## [0.1.1] - 2022-08-09

### Added

- Add the `--top` search option.

### Fixed

- Correct geographic coordinates in search results.

## [0.1.0] - 2022-08-06

### Added

- Initial release of the old tweet search CLI with CSV export, date and account filters, language and tweet type filters, and guest token rotation.

[Unreleased]: https://github.com/akiomik/squawks/compare/v0.3.1...HEAD
[0.3.1]: https://github.com/akiomik/squawks/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/akiomik/squawks/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/akiomik/squawks/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/akiomik/squawks/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/akiomik/squawks/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/akiomik/squawks/releases/tag/v0.1.0
