# Validation Entry Point Use Cases

## Validation Runner

The CAF Caller supplies a Target Project Path. The Validation Runner reads the Validation Selection from `.caf.yaml` in the Target Project root, runs each selected Validation, and returns a ValidationReport response.

```yaml
version: 1
validators:
  - oft-tracing-validator
```

The Validation Runner reports an Unknown Validation as `tool_failure` with the code `CAF-UNKNOWN-VALIDATOR` and runs the other selected Validations. Each selected Validation has one result under its name: a failed validator gives `tool_failure` with `CAF-VALIDATION-FAILED` and a message that names the `E-CAF-n` code of its error, and a result under another name gives `tool_failure` with `CAF-VALIDATION-NAME-MISMATCH`.

| Validation Name | Result |
| --- | --- |
| `oft-tracing-validator` | `pass` when the OFT Tracing Validator reports valid tracing, `finding` otherwise; `tool_failure` with `CAF-VALIDATION-FAILED` and `E-CAF-12` when OpenFastTrace cannot run |

Scenarios: [Validation Runner](../../scenarios/validation-entry-point/validation-selection/validation-runner/)

## Agent Turn Validator

A coding agent runs its Stop Hook when it is about to end its turn, and the Stop Hook calls CAF once with a Stop Event. The Agent Turn Validator runs the Validation Runner on the Target Project of the Stop Event and returns a TurnDecision response.

| Validation | TurnDecision response |
| --- | --- |
| Overall Result `pass` or `not_applicable` | `end_turn` |
| Overall Result `finding` or `tool_failure` | `continue_turn`, reason names each such result |
| Validation Runner error, such as a missing `.caf.yaml` | `continue_turn`, reason names the `E-CAF-n` code |

The reason holds the Validation Name, status, code, and message of each `finding` and `tool_failure` result. These hold no Target Project content, so the coding agent and the user can read the reason safely.

Scenarios: [Agent Turn Validator](../../scenarios/validation-entry-point/agent-turn-validation/agent-turn-validator/)

## Agent Continuation Limiter

CAF keeps no state between Stop Hook calls. It learns about earlier calls in the same Agent Turn only from the Continuation Active flag of the Stop Event. When Continuation Active is set, the Agent Continuation Limiter still runs the Agent Turn Validator and returns `end_turn` for every Validation Report and every Validation Runner error, with the same reason. This bounds the Agent Turn at one Continuation Request.

Scenarios: [Agent Continuation Limiter](../../scenarios/validation-entry-point/agent-continuation-limit/agent-continuation-limiter/)
