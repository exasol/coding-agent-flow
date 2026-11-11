# Coding Agent Flow Glossary

This chapter defines the terms used in CAF documentation and traceability.

| Term | Meaning |
| --- | --- |
| CAF | Coding Agent Flow, the language-independent workflow tool in the [CAF GitHub repository](https://github.com/exasol/coding-agent-flow). |
| OFT | OpenFastTrace, the tool CAF uses to analyze traceability links between repository artifacts. Human readers can start at the [OpenFastTrace homepage](https://openfasttrace.itsallcode.org/) and consult the [user guide](https://openfasttrace.itsallcode.org/user_guide/user_guide.html); agents can use the [OpenFastTrace GitHub repository](https://github.com/itsallcode/openfasttrace). |
| Architecture Decision Record (ADR) | A traceable record of an important decision about CAF's architecture. It states the decision question, considered alternatives, selected decision, rationale, consequences, status, and implementation evidence. An ADR uses an OFT identifier with the `adr` type, for example `adr~openfasttrace-traceability~1`. |
| Requirements specification (RS) | The set of documents that explains why CAF exists: the user or business problem it must solve, the outcomes it must provide, and the constraints it must respect. CAF stores these documents in `doc/requirements/`; they do not prescribe the technical solution. |
| Requirement | A testable, present-tense statement of CAF's current behavior or constraints. A requirement is a snapshot of the product: design explains how CAF fulfills it, implementation realizes that design, and tests provide evidence. |
| Design | Documentation of how CAF fulfills its requirements: the decisions, structures, responsibilities, and interactions chosen for the solution. Design is not the executable implementation itself. CAF stores design documentation in `doc/design/`. |
| Implementation | The executable project artifacts that realize the design. Implementation makes the designed solution available to users. |
| Test | Executable evidence that an implementation has the behavior required by a requirement or design decision. A test is distinct from the requirement: it verifies the stated intent rather than defining it. |
| Traceability | Explicit, reviewable links between requirements, design, implementation, and tests. OFT analyzes these links to reveal missing or inconsistent evidence. |
| Work Item Description (WID) | A description of an incremental unit of work, such as a GitHub issue, ticket, or epic. Work items build on one another over time. A WID provides planning context but is not automatically a requirement or a traceability artifact. |
