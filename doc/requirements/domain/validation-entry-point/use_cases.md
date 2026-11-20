# Validation Entry Point Use Cases

## Validation Runner

The CAF Caller supplies a Target Project Path. The Validation Runner reads the Validation Selection from `.caf.yaml` in the Target Project, runs each selected Validation, and returns a ValidationReport response.

```yaml
version: 1
validators:
  - oft-tracing-validator
```

The Validation Runner reports an Unknown Validation as `tool_failure` with the code `CAF-UNKNOWN-VALIDATOR` and runs the other selected Validations. Each selected Validation has one result under its name: a failed validator gives `tool_failure` with `CAF-VALIDATION-FAILED`, and a result under another name gives `tool_failure` with `CAF-VALIDATION-NAME-MISMATCH`.

| Validation Name | Result |
| --- | --- |
| `oft-tracing-validator` | `pass` when the OFT Tracing Validator reports valid tracing, `finding` otherwise |

Scenarios: [Validation Runner](../../scenarios/validation-entry-point/validation-selection/validation-runner/)
