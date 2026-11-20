# OFT Validation Glossary

Shared CAF terms are defined in the [Coding Agent Flow Domain Glossary](../glossary.md).

| Term | Meaning |
| --- | --- |
| Traceability | Reviewable links between OFT Specification Items. |
| OFT Specification Item | An OpenFastTrace item with a type, identifier, and revision number. |
| OFT Feature | An OFT Specification Item with the `feat` type that describes a product capability. |
| OFT Requirement | An OFT Specification Item with the `req` type that describes a product requirement and covers an OFT Feature. |
| OFT Artifact Type | The type of an OFT Specification Item, such as `req` or `feat`. |
| OFT Specification Item ID | The identifier of one OFT Specification Item. |
| OFT Requirements Detector | The use case that identifies OFT Requirements from the `req` specification item IDs in an OFT Validation response. |
| OFT Tracing Validator | The use case that validates OFT tracing in a Target Project. |
| OFT Validation response | The immutable result of validating a Target Project with OFT. `valid` reports the validation result. `oftSpecificationItemIDs` maps each OFT Artifact Type to its OFT Specification Item IDs. |
| OFT Trace | The OFT Specification Items of a Target Project with their OFT Coverage, and the OFT Result, as OFT reports them. |
| OFT Coverage | The OFT Artifact Types an OFT Specification Item needs, and the OFT Specification Items it covers. |
| OFT Result | `success` when OFT finds no tracing defect, `failure` otherwise. |
