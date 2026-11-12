# Spec Validation Glossary

Shared CAF and OFT terms are defined in the
[Coding Agent Flow Domain Glossary](../glossary.md).

| Term | Meaning |
| --- | --- |
| Requirement Document | The domain model for one discovered OFT item whose identifier type is `req`. It owns the non-blank Needs and Rationale invariants. |
| Requirement Document Structure Validator | The use case that evaluates each Requirement Document in a supplied Source Project and returns one RequirementStructure response. |
| Spec Structure | The set of structural invariants evaluated for Requirement Documents in a Source Project: each has non-blank Needs and Rationale fields. |
| RequirementStructure response | The immutable response to one Requirement Document Structure Validator invocation. It has a boolean `valid` field that is true only if every discovered Requirement Document satisfies the Spec Structure invariants. |
