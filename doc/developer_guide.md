# Developer Guide

## Machine Validation CLI

`doc~caf-machine-validation~1`

Covers:

- `req~caf-machine-validation~1`

CAF provides one non-interactive validation command:

```shell
caf validate
```

Without an agent option, the command is intended for CI. It prints a JSON report and exits with 0 for pass, 1 for a validation finding, or 2 for a tool failure.

## Project Configuration Discovery

`doc~caf-project-discovery~1`

Covers:

- `req~caf-project-discovery~1`
- `req~caf-validation-routing~1`

CAF searches from the invocation working directory toward the filesystem root for `.caf.yaml`. The nearest file defines the project root and validation policy:

```yaml
version: 1
validators:
  - requirements-document
```

The initial supported validator identifier is `requirements-document`. OFT execution will be added by the OFT work package without adding another public command.

## Codex Stop Hook

`doc~caf-codex-stop-hook~1`

Covers:

- `req~caf-agent-selection~1`
- `req~caf-codex-stop-response~1`

The initial agent integration supports Codex:

```shell
caf validate --agent codex
```

The Codex `Stop` hook supplies its JSON object through standard input. CAF uses the payload's `cwd` for project discovery and writes only Codex response JSON to standard output. A pass produces `{}`. A finding or tool failure produces a `decision: block` response so Codex continues working. If `stop_hook_active` is already true, CAF allows the turn to stop to prevent an uncontrolled continuation loop.

Hook installation is outside this command's scope. See the [Codex Hooks documentation](https://learn.chatgpt.com/docs/hooks) for the current provider protocol.

## Requirements Document Detection

CAF recognizes a direct regular file in `doc/requirements/` as a requirements
document when it uses one of these extensions:

- `.md` for Markdown;
- `.rst` for reStructuredText.

The check does not search subdirectories or validate document content.

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
