# Coding Agent Flow Glossary

This chapter defines the terms used in CAF documentation and traceability.

| Term | Meaning |
| --- | --- |
| CAF | Coding Agent Flow, the language-independent workflow tool in the [CAF GitHub repository](https://github.com/exasol/coding-agent-flow). |
| OFT | OpenFastTrace, the tool CAF uses to analyze traceability links between repository artifacts. See the [OpenFastTrace GitHub repository](https://github.com/itsallcode/openfasttrace). |
| Requirements specification (RS) | The set of documents that explains why CAF exists: the user or business problem it must solve, the outcomes it must provide, and the constraints it must respect. CAF stores these documents in `doc/requirements/`; they do not prescribe the technical solution. |
| Requirement | One testable statement of a reason CAF is needed, an outcome it must provide, or a constraint it must meet. A requirement starts a trace: design explains how CAF will fulfill it, implementation realizes that design, and tests provide evidence. |
| Design | Documentation of how CAF fulfills its requirements: the decisions, structures, responsibilities, and interactions chosen for the solution. Design is not the executable implementation itself. CAF stores design documentation in `doc/design/`. |
| Implementation | The executable project artifacts that realize the design, principally CAF's Go source code. Implementation makes the designed solution available to users. |
| Test | Executable evidence that an implementation has the behavior required by a requirement or design decision. A test is distinct from the requirement: it verifies the stated intent rather than defining it. |
| Traceability | Explicit, reviewable links between requirements, design, implementation, and tests. OFT analyzes these links to reveal missing or inconsistent evidence. |
| Work Item Description (WID) | A description of a unit of work, such as a GitHub issue, ticket, or epic. A WID provides planning context but is not automatically a requirement or a traceability artifact. |
