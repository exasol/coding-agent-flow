# Spec Validation Bounded Context

Spec Validation is a read-only bounded context that checks the structure of
Requirement Documents discovered in a supplied Source Project.

| Document | Purpose |
| --- | --- |
| [glossary.md](glossary.md) | Defines vocabulary specific to Spec Validation. |
| [context_and_scope.md](context_and_scope.md) | Defines the boundary, actor, and scope. |
| [entity_model_relationship.md](entity_model_relationship.md) | Defines entities, relationships, and the RequirementStructure response model. |
| [value_objects.md](value_objects.md) | Defines immutable domain concepts and validation boundaries. |
| [invariants.md](invariants.md) | Defines Requirement Document structural rules and result calculation. |
| [use_cases.md](use_cases.md) | Defines the public use case and links to its acceptance scenarios. |
| [failures.md](failures.md) | Records intentionally deferred failure-contract decisions. |

Shared vocabulary, including the meaning of a Spec, is defined in the
[Coding Agent Flow Domain Glossary](../glossary.md). The feature and requirement are defined
in [`../../system_requirements.md`](../../system_requirements.md).
