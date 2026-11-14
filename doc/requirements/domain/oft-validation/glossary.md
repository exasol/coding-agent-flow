# OFT Validation Glossary

Shared CAF and OFT terms are defined in the [Coding Agent Flow Domain Glossary](../glossary.md).

| Term | Meaning |
| --- | --- |
| Traceability | Reviewable links between OFT elements. |
| OFT Element | An OpenFastTrace item with a type and identifier. |
| OFT Feature | An OFT Element with the `feat` type that describes a product capability. |
| OFT Requirement | An OFT Element with the `req` type that describes a product requirement and covers an OFT Feature. |
| OFT Element Type | The type of an OFT Element, such as `req` or `feat`. |
| OFT Element ID | The identifier of one OFT Element. |
| OFT Requirements Detector | The use case that identifies OFT Requirements from the `req` element IDs in an OFTValidation response. |
| OFT Tracing Validator | The use case that validates OFT tracing in a Target Project. |
| OFTValidation response | The immutable result of validating a Target Project with OFT. `valid` reports the validation result. `oftElementIDs` maps each OFT Element Type to its OFT Element IDs. |
