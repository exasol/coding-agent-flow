# CAF Requirements

## OFT Requirements Detection

`req~oft-requirements-detection~1`

CAF uses OFT to recognize OFT Requirements in a Target Project.

Needs: scn, impl, utest

Rationale: CAF uses a consistent input set for OFT validation.

## OFT Tracing Validation

`req~oft-tracing-validation~1`

CAF validates OFT Specification Items in the Target Project and returns an `OFTValidation` response. `valid` reports the validation result. `oftSpecificationItemIDs` maps each OFT Artifact Type to its OFT Specification Item IDs.

Covers:

- `feat~oft-validation~1`

Needs: scn, impl

Rationale: CAF reports traceability results before implementation begins.

## Link Validation

`req~link-validation~1`

CAF checks the Links in the `.md`, `.rst`, `.txt`, and `.xml` documents of a Target Project and returns a `LinkValidation` response. `valid` reports the result. `brokenLinks` lists each Broken Link with source file, line, target, and status. CAF leaves the Target Project unchanged and ignores link-checker configuration files in it.

Covers:

- `feat~code-validation~1`

Needs: dsn, scn, impl, utest

Rationale: CAF reports broken document links before an agent hands over its work.

## External Link Validation

`req~external-link-validation~1`

CAF checks Local Links offline by default. On request, CAF also checks External Links and counts a URL that answers `429 Too Many Requests` as valid.

Covers:

- `feat~code-validation~1`

Needs: dsn, scn, impl, utest

Rationale: The default check is fast and independent of the network, and the caller adds External Links when needed.

## Validation Selection

`req~validation-selection~1`

CAF runs each validation that the Target Project selects and returns a `ValidationReport` with one result per selected validation, in the selected order. CAF reports a tool failure for an unknown validation.

Covers:

- `feat~validation-entry-point~1`

Needs: scn, impl, utest

Rationale: Each Target Project chooses its validations, and new validations use the same entry point.

## Validation Result

`req~validation-result~1`

CAF reports each validation as `pass`, `finding`, `tool_failure`, or `not_applicable`, and the overall result as the most severe one, in the order `tool_failure`, `finding`, `pass`, `not_applicable`.

Covers:

- `feat~validation-entry-point~1`

Needs: scn, impl, utest

Rationale: Coding agents and users act on one clear result for all selected validations.
