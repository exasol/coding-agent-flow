# Developer Guide

## Requirements Document Detection

CAF recognizes a direct regular file in `doc/requirements/` as a requirements
document when it uses one of these extensions:

- `.md` for Markdown;
- `.rst` for reStructuredText.

The check does not search subdirectories or validate document content.

## Code Validation Code Structure

The Link Checker lives in `internal/codevalidation/`:

| Package | Content | Example |
|---|---|---|
| `domain` | Domain classes and their rules | `LinkValidation`, `BrokenLink`, `LinkScope` |
| `application/port` | Interfaces the use cases call | `BrokenLinkFinder` |
| `application/service` | One service per use case | `LinkChecker` |
| `adapter/lychee` | Implementation of `BrokenLinkFinder` with lychee | `Adapter` |

The root package `codingagentflow` is the public API. Its constructor passes the adapter and the Link Scope to the use case:

```go
func NewLinkChecker(options ...LinkCheckerOption) LinkChecker {
	settings := linkCheckerSettings{scope: domain.LocalLinks}
	for _, option := range options {
		option(&settings)
	}
	return service.NewLinkChecker(lychee.NewAdapter(), settings.scope)
}
```

`internal/codevalidation/architecture_test.go` checks the imports of each layer. `internal/global/target_project.go` holds the `TargetProject` that all domains share.

### lychee Command

The adapter runs lychee in the Target Project:

```shell
lychee --config <CAF configuration> --format json --no-progress --include-fragments --extensions md,rst,txt,xml --root-dir <target project> <scope arguments> .
```

| Link Scope | Scope arguments |
|---|---|
| Local Links: `NewLinkChecker()` | `--offline` |
| Local and External Links: `NewLinkChecker(WithExternalLinks())` | `--accept 200..=299,429 --timeout 10 --max-retries 2` |

The CAF configuration is a temporary file outside the Target Project. It replaces the `lychee.toml` that lychee reads from its working directory by default.

lychee uses `GITHUB_TOKEN` from the environment for GitHub URLs, which raises the GitHub rate limit.

The adapter reads the Broken Links from the JSON report for exit code `0` (all Links resolve) and `2` (Broken Links). It reports other exit codes as errors.

### lychee Setup

CAF uses the lychee binary from `LYCHEE_BINARY`. If the variable is unset, CAF uses lychee 0.24.2 from the user cache directory, for example `~/.cache/coding-agent-flow/lychee/0.24.2/lychee` on Linux. On first use, CAF downloads it for Linux or macOS and verifies its SHA-256 checksum.

Before each check, CAF requires `lychee --version` to report `lychee 0.24.2`, the version whose JSON report the adapter reads. Otherwise, it fails with `E-CAF-4`.

Run the acceptance tests of the Link Checker:

```shell
go test ./test/acceptance/code-validation/...
```

The acceptance tests check URLs against a local web server on `127.0.0.1` and run offline once lychee is in the cache.

## Continuous Integration

The GitHub Actions workflow in `.github/workflows/ci-build.yml` runs on every
pull request and on pushes to `main`. It uses Go 1.26.8 and verifies:

* all packages build; and
* unit tests pass.

The repository's root `go.mod` defines its module. Add `go.sum` when the
project introduces external dependencies.

## Managing Release Notes

Keep one Markdown file per released version under `doc/changes/`, named
`changes_<version>.md`, for example `changes_0.2.0.md`. Add the newest version
as the first entry in `doc/changes/changelog.md`.

Prepare the next release-note file in the release preparation pull request.
Its heading contains the project name, semantic version, and release date; it
should have a short summary followed by only the relevant sections, such as
`Features`, `Bugfixes`, `Security`, `Documentation`, and `Dependency Updates`.
Reference related GitHub issues or pull requests in the entries. Do not rewrite
notes for already released versions; publish corrections as a new release if
they affect users.

Use semantic versioning for the public Go module API:

* patch: backward-compatible bug fixes, dependency or security updates;
* minor: backward-compatible features; and
* major: incompatible exported-API changes. Go modules require a `/vN` module
  path suffix for major versions 2 and later.

## Releasing

Create GitHub releases with
[Release Droid](https://github.com/exasol/release-droid), using a generic
language configuration. Go users consume the module by its Git tag, so there
is no binary artifact to upload unless the project later introduces a CLI
distribution.

Before each release:

1. Choose the semantic version and update the release-note file and changelog.
2. Update every project-owned version source to that version. When a runtime
   version is introduced, keep it in `internal/version/version.go`; configure
   the same version in `.project-keeper.yml`.
3. Ensure the CI workflow is green and run the release validation locally:

   ```shell
   java -jar .../release-droid-*.jar -n coding-agent-flow --local . --goal validate
   ```

4. Merge the prepared release into `main`, then create the GitHub release:

   ```shell
   java -jar .../release-droid-*.jar -n coding-agent-flow --goal release --skipvalidation
   ```

5. Verify that the created tag is the intended Go module version (for example,
   `v0.1.0`) and that the GitHub release text links to the matching change file.

`release_config.yml` limits Release Droid to GitHub and the generic release
mode. It intentionally does not publish assets.
