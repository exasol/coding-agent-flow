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
| Agent Turn | One unit of work of a coding agent, ending when the agent stops. CAF sees an Agent Turn only through Stop Events. |
| Stop Hook | The command a coding agent runs when it is about to stop. It calls CAF with a Stop Event. |
| Stop Event | What CAF receives from a Stop Hook: the Target Project Path and the Continuation Active flag. |
| Continuation Active | The Agent Turn already continued because of a Stop Hook (Codex: `stop_hook_active`). |
| Continuation Request | CAF's request that the agent continues its turn, with the reason. |
| Continuation Limit | CAF sends at most one Continuation Request per Agent Turn: it lets the agent stop when Continuation Active is set. |
| Turn Decision | `end_turn`, or `continue_turn` with a reason. |
| Agent Turn Validator | The use case that validates the Target Project of a Stop Event and returns a Turn Decision. |
| Agent Continuation Limiter | The use case that applies the Continuation Limit: it lets the agent end its turn when Continuation Active is set. |
| TurnDecision response | The immutable result of the Agent Turn Validator and the Agent Continuation Limiter: `action` (`end_turn`, `continue_turn`), `reason`, and the Validation Report. |
