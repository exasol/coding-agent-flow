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
