# System Requirements

## Requirements Document Detection

`req~requirements-document-detection~1`

CAF shall recognize a direct regular file in `doc/requirements/` as a requirements document when its extension is supported.

Supported extensions:

- `.md`
- `.rst`

Needs: impl, utest

## Machine Validation Entry Point

`req~caf-machine-validation~1`

CAF shall provide `caf validate` as its single non-interactive validation entry point.

Rationale: Agent hooks and CI need one stable command that can run every validation selected by project policy.

Needs: dsn, impl, itest, doc

## Project Configuration Discovery

`req~caf-project-discovery~1`

CAF shall derive the project root from the nearest `.caf.yaml` found by searching from the invocation working directory toward the filesystem root.

Rationale: Hook commands must remain independent of checkout-specific project-root paths.

Needs: dsn, impl, utest, itest, doc

## Codex Agent Selection

`req~caf-agent-selection~1`

CAF shall select the Codex Stop-hook protocol when `caf validate` receives `--agent codex`.

Rationale: Explicit selection is deterministic and does not depend on guessing the provider from overlapping JSON fields.

Needs: dsn, impl, itest, doc

## Codex Stop Response

`req~caf-codex-stop-response~1`

CAF shall render the validation decision as a valid Codex Stop-hook JSON response.

Rationale: Codex must be able to consume a pass or continuation decision without parsing human-oriented output.

Needs: dsn, impl, itest, doc

## Configured Validation Routing

`req~caf-validation-routing~1`

CAF shall run every supported validator declared in the discovered project configuration.

Rationale: Adding a validator must not require another public hook command.

Needs: dsn, impl, utest, itest, doc
