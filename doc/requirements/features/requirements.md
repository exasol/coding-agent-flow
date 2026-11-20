# CAF Requirements

## OFT Requirements Detection

`req~oft-requirements-detection~1`

CAF uses OFT to recognize OFT Requirements in a Target Project.

Needs: scn, impl, utest

Rationale: CAF uses a consistent input set for OFT validation.

## OFT Tracing Validation

`req~oft-tracing-validation~1`

CAF validates OFT Specification Items in the Target Project and returns an `OFTValidation` response. `valid` reports the validation result. `oftSpecificationItemIDs` maps each OFT Artifact Type to its OFT Specification Item IDs. CAF returns an error when it can't trace the Target Project, resolves a symbolic link as Target Project path, and leaves the Target Project unchanged.

Covers:

- `feat~oft-validation~1`

Needs: scn, impl, utest

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
