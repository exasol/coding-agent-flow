# Spec Validation Value Objects

Value objects are immutable domain concepts defined by their values rather
than an independent identity or lifecycle. They are owned by the entity or use
case named below.

| Value object | Owner | Value | Validation boundary |
| --- | --- | --- | --- |
| Source Project Path | Source Project | The caller-supplied path identifying a Source Project for one validation request. | Path normalization and filesystem accessibility are outside this contract. |
| OFT Requirement Identifier | Requirement Document | The identifier of the discovered OFT item. | Its OFT item type is `req`. |
| Needs | Requirement Document | The `Needs` field value. | It is non-blank after trimming surrounding whitespace. Trace-target syntax is outside this contract. |
| Rationale | Requirement Document | The `Rationale` field value. | It is non-blank after trimming surrounding whitespace. Meaning beyond non-blankness is outside this contract. |

Requirement Document owns the Needs and Rationale value objects. A missing,
empty, or whitespace-only value violates that Requirement Document's
structural invariants.
