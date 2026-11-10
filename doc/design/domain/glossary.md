# Coding Agent Flow Glossary

This chapter defines the terms used in CAF documentation and traceability.

| Term | Meaning |
| --- | --- |
| CAF | Coding Agent Flow, the language-independent workflow tool in the [CAF GitHub repository](https://github.com/exasol/coding-agent-flow). |
| OFT | OpenFastTrace, the tool CAF uses to analyze traceability links between repository artifacts. See the [OpenFastTrace GitHub repository](https://github.com/itsallcode/openfasttrace). |
| Requirements specification (RS) | The set of documents that state what CAF must do and the constraints it must meet. CAF stores these documents in `doc/requirements/`; a requirements document expresses intended, externally meaningful behavior rather than implementation details. |
| Requirement | One testable statement of a needed capability, behavior, or constraint. A requirement is the starting point for a trace: design explains how it is addressed, implementation realizes it, and tests provide evidence. |
| Design | Documentation of the decisions, structures, responsibilities, and interactions chosen to satisfy requirements. Design explains *how and why* CAF is organized; it is not the implementation itself. CAF stores design documentation in `doc/design/`. |
| Implementation | The executable project artifacts that realize a design, principally CAF's Go source code. Implementation turns the intended behavior described by requirements and design into behavior users can run. |
| Test | Executable evidence that an implementation has the behavior required by a requirement or design decision. A test is distinct from the requirement: it verifies the stated intent rather than defining it. |
| Traceability | Explicit, reviewable links between requirements, design, implementation, and tests. OFT analyzes these links to reveal missing or inconsistent evidence. |
| Work Item Description (WID) | A description of a unit of work, such as a GitHub issue, ticket, or epic. A WID provides planning context but is not automatically a requirement or a traceability artifact. |
