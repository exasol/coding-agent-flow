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
