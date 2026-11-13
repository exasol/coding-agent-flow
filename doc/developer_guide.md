# Developer Guide

## Reconcile Diagrams

Use the `reconcile diagrams` command to apply the project-local
`$reconcile-diagrams` skill when a CAF domain diagram must be added, repaired,
or reviewed. The skill keeps Mermaid diagrams in their owning bounded-context
documents under `doc/requirements/domain/<bounded-context>/`:

- `context_and_scope.md` answers who interacts with the context and what
  crosses its boundary;
- `entity_model_relationship.md` answers how entities, value objects, and
  relationships fit together; and
- `use_cases.md` answers how an actor invokes a use case and receives its
  outcome.

The command reconciles a diagram with the glossary, model, invariants, use case,
and Gherkin scenarios, and synchronizes each diagram with its corresponding
bounded-context documents. It does not infer undocumented technical behavior
or replace the surrounding prose and specifications. CAF currently has no
pinned Mermaid renderer; do not download one merely to validate a diagram.
Review the rendered Markdown in the normal documentation viewer and report
that renderer gap when it applies.

## Reconcile Domain Model

Use the `reconcile domain model` command to apply the project-local
`$reconcile-domain-model` skill when Gherkin scenarios and domain
documentation must be synchronized. The command resolves a scenario list's
bounded context from
`doc/requirements/scenarios/<bounded-context>/<feature>/<use-case>/` and
reconciles it with `doc/requirements/domain/<bounded-context>/`. When several
scenario lists are selected, it groups them by bounded context and reconciles
each context independently. It keeps global vocabulary in
`doc/requirements/domain/glossary.md`, context vocabulary in the matching
context glossary, and canonical Gherkin only in `doc/requirements/scenarios/`.

## Reconcile Acceptance Tests

Use the `reconcile acceptance tests` command to apply the project-local
`$reconcile-acceptance-tests` skill when acceptance evidence must be aligned
with Gherkin scenarios, use cases, and the domain model. The skill resolves
the bounded context from the scenario path, maps its scenarios to one
black-box Go acceptance test each, and checks that tests use an approved public
boundary rather than an internal package. It creates those tests before
production code and fixtures, so new-feature and bug-fix tests are expected to
be red until their public binding exists. The generic runner reports missing
bindings or fixtures as executable failures, rather than compilation errors or
skipped tests. It keeps `.feature` files as the canonical BDD specification.
CAF's repository-owned Go runner is for CAF's own public use cases: it maps a
canonical scenario ID to a fixture, one public-boundary invocation, and its
assertions. It is not a cross-language test-execution framework. The runner
must preserve canonical scenario IDs and requires deterministic binding
validation.

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
