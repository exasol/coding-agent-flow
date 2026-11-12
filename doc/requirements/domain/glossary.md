# Coding Agent Flow Domain Glossary

This glossary defines terms shared across CAF bounded contexts. Each bounded
context owns its domain-specific vocabulary in its own `glossary.md` file.

| Term | Meaning |
| --- | --- |
| Bounded context | A cohesive area of the domain that owns its vocabulary, model, rules, and use cases. It exchanges explicit contracts rather than live domain objects with other contexts. |
| Domain model | The human-readable representation of a bounded context's concepts, relationships, invariants, and outcomes. |
| Entity | A domain concept with an identity that remains meaningful even when its attributes change. |
| Value object | An immutable domain concept defined by its value rather than an independent identity or lifecycle. |
| Source Project | A project selected as input to a CAF operation. CAF may read it without owning or modifying it. |
| Source Project Path | An immutable value object holding the caller-supplied path that identifies a Source Project for one operation. It does not normalize the path or establish filesystem accessibility. |
| Requirements document | A source artifact in a Source Project that can declare OFT items. Recognition policy is owned by the context that consumes it. |
| Spec | Structured specification documentation that records a product capability, its requirements, and the evidence expected to realize and verify them. In OFT tracing, a Spec connects a `feat` item to its `req` items; a requirement's `Needs` field declares expected downstream evidence such as `scn`, `impl`, and `at`; scenarios, implementation, and acceptance tests then provide traceable coverage back to that requirement. A Spec is the human-maintained source of intent and traceability, not implementation or test evidence itself. |
| OFT Requirement | An OFT item whose identifier type is `req`, such as `req~spec-structure-validation~1`. |
| OFT Requirement Identifier | An immutable value object identifying an OFT Requirement. Its OFT item type is `req`. |
| Needs | The OFT field that states downstream artifact types expected to provide traceability evidence for an item. A bounded context may define additional syntax or non-blankness rules. |
| Rationale | The OFT field that records why an item is needed. A bounded context may define additional content or non-blankness rules. |
| Non-blank | A text value containing at least one non-whitespace character after surrounding whitespace is removed. |
