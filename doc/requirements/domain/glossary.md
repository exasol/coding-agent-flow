# Spec Validation Glossary

| Term | Meaning |
| --- | --- |
| Source Project | The project identified by the source-project path supplied for one validation request. CAF reads it but does not own or modify it. |
| Requirements document | A recognized source artifact in a Source Project that can declare OFT items. |
| Requirement Document | The domain model for one discovered OFT item whose identifier type is `req`. It owns the non-blank Needs and Rationale invariants. |
| Spec | Structured specification documentation that records a product capability, its requirements, and the evidence expected to realize and verify them. In OFT tracing, a Spec connects a `feat` item to its `req` items; a requirement's `Needs` field declares the expected downstream evidence, such as `scn`, `impl`, and `at`; scenarios, implementation, and acceptance tests then provide traceable coverage back to that requirement. A Spec is the human-maintained source of intent and traceability, not the implementation or the test evidence itself. |
| OFT Requirement | An OFT item whose identifier type is `req`, such as `req~requirement-structure-validation~1`; it is represented by a Requirement Document in the Spec Validation domain. |
| Needs | The OFT field that states the downstream artifact types required to provide traceability evidence for an item. For an OFT Requirement, it must be non-empty after trimming surrounding whitespace. |
| Rationale | The OFT field that records why an item is needed. For an OFT Requirement, it must be non-empty after trimming surrounding whitespace. |
| Non-blank | A text value containing at least one non-whitespace character after surrounding whitespace is removed. |
| Spec Structure Validator | The use case that evaluates each Requirement Document in a supplied Source Project and returns one RequirementStructure response. |
| Spec Structure | The set of structural invariants evaluated for Requirement Documents in a Source Project: each has non-blank Needs and Rationale fields. |
| RequirementStructure response | The immutable response to one Spec Structure Validator invocation. It has a boolean `valid` field that is true only if every discovered Requirement Document satisfies the Spec Structure invariants. |
