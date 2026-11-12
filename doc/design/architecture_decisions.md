# Architecture Decisions

This chapter records the decisions that determine how Coding Agent Flow is structured and evolves.

Record a decision when it is hard to change later, likely to be debated again, relevant for cost or risk, or necessary to explain a non-obvious implementation choice.

## Traceability

### How does Coding Agent Flow keep its artifacts traceable?

Coding Agent Flow uses [OpenFastTrace](https://github.com/itsallcode/openfasttrace) to analyze traceability links between system requirements, architecture decisions, implementation, and tests.

This decision applies to Coding Agent Flow's own artifacts. Its system requirements record what Coding Agent Flow must do, and its architecture decisions record how it will do it. Its implementation is source code, while its tests provide executable evidence. Contributors link these artifacts using stable OFT identifiers and coverage relationships. OFT analyzes those links so that a change can be followed from intent to evidence and gaps can be detected.

Coding Agent Flow keeps system requirements in `doc/requirements/` and architecture decisions in `doc/design/`. Source code and tests remain alongside the packages they implement and verify. Each traceable artifact uses stable OpenFastTrace identifiers and links only to the artifacts it actually fulfills or verifies.

Considered alternatives:

1. Use OpenFastTrace to analyze explicitly declared traceability links.
2. Do not use a traceability analysis tool.

#### Traceability with OpenFastTrace

`adr~openfasttrace-traceability~1`

The system uses OpenFastTrace to analyze explicitly declared traceability links.

### Rationale

Coding Agent Flow is a workflow tool whose behavior will grow across requirements detection, agent orchestration, and related project checks. Traceability makes the reasoning behind that behavior reviewable and helps prevent system requirements, architecture decisions, code, and tests from drifting apart.

### Consequences

Changes to the architecture should update the affected architecture decisions before or alongside the code and tests. Contributors must preserve existing OpenFastTrace identifiers when their meaning is unchanged and add links only where they have supporting evidence.

Needs: impl

Tags: `traceability`, `openfasttrace`

## Open Issues

* None.
