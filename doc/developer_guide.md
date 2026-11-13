# Developer Guide

## OFT Requirements Detection

CAF passes the Target Project path to OFT. OFT recognizes OFT requirement IDs in Markdown (`.md`) and reStructuredText (`.rst`) files and returns them in the `req` entry of `OFTValidation.oftElementIDs`.

OFT does not return IDs from plain-text (`.txt`) files.

## Continuous Integration

The GitHub Actions workflow in `.github/workflows/ci-build.yml` runs on every pull request and on pushes to `main`. It uses Go 1.26.8 and verifies:

* all packages build; and
* unit and acceptance tests pass.

The repository's root `go.mod` defines its module. Add `go.sum` when the project introduces external dependencies.

## Acceptance Tests

The generic acceptance runner in `test/acceptance/acceptance_test.go` reads canonical Gherkin scenarios and executes public-boundary programs. Fixture data supplies temporary Target Project files. Typed bindings define the expected public response. The runner does not import planned APIs, so scenarios compile before implementation. A missing boundary or fixture is an expected-red result with a clear failure message.

Each acceptance binding carries an `at` marker that covers its Gherkin scenario. Acceptance bindings create Target Project fixtures in temporary directories.

## Managing Release Notes

Keep one Markdown file per released version under `doc/changes/`, named `changes_<version>.md`, for example `changes_0.2.0.md`. Add the newest version as the first entry in `doc/changes/changelog.md`.

Prepare the next release-note file in the release preparation pull request. Include the project name, semantic version, release date, short summary, relevant sections, and related issues or pull requests. Publish corrections as a new release when they affect users.

Use semantic versioning for the public Go module API:

* patch: backward-compatible bug fixes, dependency or security updates;
* minor: backward-compatible features; and
* major: incompatible exported-API changes. Go modules require a `/vN` module path suffix for major versions 2 and later.

## Releasing

Create GitHub releases with [Release Droid](https://github.com/exasol/release-droid) and the generic language configuration. Go users consume the module by its Git tag.

Before each release:

1. Choose the semantic version and update the release-note file and changelog.
2. Update every project-owned version source to that version. When a runtime version is introduced, keep it in `internal/version/version.go`; configure the same version in `.project-keeper.yml`.
3. Ensure the CI workflow is green and run the release validation locally:

   ```shell
   java -jar .../release-droid-*.jar -n coding-agent-flow --local . --goal validate
   ```

4. Merge the prepared release into `main`, then create the GitHub release:

   ```shell
   java -jar .../release-droid-*.jar -n coding-agent-flow --goal release --skipvalidation
   ```

5. Verify that the created tag is the intended Go module version (for example, `v0.1.0`) and that the GitHub release text links to the matching change file.

`release_config.yml` configures GitHub and the generic release mode.
