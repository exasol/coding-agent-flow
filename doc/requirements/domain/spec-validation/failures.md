# Spec Validation Failure Contract

## Status

Failure handling is intentionally deferred. The current domain contract only
defines a RequirementStructure result for structurally evaluable Requirement
Documents.

## Deferred conditions

The following conditions require a later CLI and runtime contract before they
become observable outcomes:

| Condition | Deferred decision |
| --- | --- |
| Source Project Path is unavailable | Diagnostic representation and process-exit behavior. |
| Requirements document cannot be discovered or read | Whether this is a finding or a tool failure. |
| Document changes while it is read | Consistency and retry behavior. |
| OFT item cannot be parsed | Classification, diagnostic code, and response shape. |

## Future failure-contract template

For every newly supported failure, record the triggering condition, owner,
stable diagnostic code, safe message, machine response, process-exit behavior,
and acceptance scenario. Do not infer these behaviors from the structural
`valid` result.
