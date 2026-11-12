# Spec Validation Invariants

This document is the rule catalog for the Requirement Document Structure
Validator use case. The rules implement
`req~requirement-structure-validation~1`.

| Rule | Owner | Statement |
| --- | --- | --- |
| Requirement document identity | Requirement Document | Its OFT Requirement Identifier has type `req`. |
| Needs is non-blank | Requirement Document | Its Needs value contains at least one non-whitespace character after trimming surrounding whitespace. |
| Rationale is non-blank | Requirement Document | Its Rationale value contains at least one non-whitespace character after trimming surrounding whitespace. |
| Result calculation | RequirementStructure response | `valid` is `true` only when every discovered Requirement Document satisfies the preceding rules. |
| Empty result | RequirementStructure response | `valid` is `true` when no Requirement Documents are discovered. |

The rule catalog does not validate `Needs` trace-target syntax or the semantic
content of a Rationale. See [value_objects.md](value_objects.md) for the
owned values and [use_cases.md](use_cases.md) for behavioral scenarios.
