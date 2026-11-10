# Architectural Requirements

This chapter records architecture decisions and requirements that guide CAF's
structure and evolution.

## ADR: Traceability with OpenFastTrace

### Decision

CAF uses [OpenFastTrace (OFT)](https://github.com/itsallcode/openfasttrace) to
analyze traceability between its requirements, design, implementation, and
tests.

Requirements are the source of intended behavior and constraints. Design
documents explain the architectural decisions that fulfill them. Implementation
is the Go source code that realizes those decisions, while tests are executable
evidence that the implementation exhibits the required behavior. OFT links
these artifacts so that a change can be followed from intent to evidence and
gaps can be detected.

CAF keeps requirements specifications in `doc/requirements/` and design
documentation in `doc/design/`. Source code and Go tests remain alongside the
packages they implement and verify. Each traceable artifact uses stable OFT
identifiers and links only to the artifacts it actually fulfills or verifies.

### Rationale

CAF is a workflow tool whose behavior will grow across requirements detection,
agent orchestration, and related project checks. Traceability makes the
reasoning behind that behavior reviewable and helps prevent requirements,
design, code, and tests from drifting apart.

### Consequences

Changes to behavior should update the affected requirements and design before
or alongside the code and tests. Contributors must preserve existing OFT IDs
when their meaning is unchanged and add links only where they have supporting
evidence.

### Status

Accepted.
