# Validation Entry Point Glossary

Shared CAF terms are defined in the [Coding Agent Flow Domain Glossary](../glossary.md).

| Term | Meaning |
| --- | --- |
| Validation Entry Point | The single way coding agents and users run the validations a Target Project selects. |
| Validation | A check that CAF runs on a Target Project, such as `oft-tracing-validator`, identified by its Validation Name. |
| Validation Name | The name of a Validation in `.caf.yaml`. |
| Validation Selection | The Validation Names that a Target Project lists in `.caf.yaml`, in order. |
| Unknown Validation | A selected Validation Name that CAF has no Validation for. |
| Validation Result Status | `pass`: the Validation completed and accepts the Target Project; `finding`: the Validation completed and found a problem; `tool_failure`: the Validation could not complete reliably; `not_applicable`: the Validation applies to other Target Projects. |
| Validation Result | The result of one Validation: Validation Name, Validation Result Status, code, and message. |
| Overall Result | The most severe Validation Result Status of a Validation Report, in the order `tool_failure`, `finding`, `pass`, `not_applicable`. |
| Validation Runner | The use case that runs the Validation Selection of a Target Project and returns a Validation Report. |
| ValidationReport response | The immutable result of the Validation Runner: `status`, the Overall Result, and `results`, one Validation Result per selected Validation in the selected order. |
