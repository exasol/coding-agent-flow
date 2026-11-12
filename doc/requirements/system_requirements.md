# System Requirements

## Requirements Document Detection

`req~requirements-document-detection~1`

CAF recognizes a direct regular file in `doc/requirements/` as a requirements document when its extension is supported.

Supported extensions:

- `.md`
- `.rst`

Needs: impl, utest

Rationale: CAF needs a predictable way to identify requirements documents
before it can validate their contents.

## Spec Validation

`feat~spec-validation~1`

CAF validates declared OFT requirement metadata in requirements documents
for a supplied source project and returns a machine-facing validation result.

Status: draft

Needs: req

### Spec Structure Validation

`req~spec-structure-validation~1`

CAF inspects every OFT requirement item (`req~...~revision`) in every
recognized requirements document of a supplied source project. For each item,
CAF requires non-empty `Needs` and `Rationale` fields. CAF returns a
`RequirementStructure` response whose `valid` boolean field is `true` only
when every inspected OFT requirement satisfies both field rules.

Status: draft

Covers:

- `feat~spec-validation~1`

Needs: scn, impl, at

Rationale: Requirements need both an explicit downstream traceability
expectation and a recorded reason so that automated checks can identify
incomplete requirement records before implementation begins.
