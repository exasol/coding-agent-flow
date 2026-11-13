# Machine Validation CLI Design

`dsn~caf-validation-cli~1`

Covers:

- `req~caf-machine-validation~1`
- `req~caf-project-discovery~1`
- `req~caf-agent-selection~1`
- `req~caf-codex-stop-response~1`
- `req~caf-validation-routing~1`

## Scope

The initial machine interface supports CI and the Codex `Stop` hook. Shared events, results, decisions, and validators remain independent of Codex so another provider can be added later without changing validation logic. This work does not install hooks or implement OFT execution.

## Command Contract

| Invocation | Caller | Input | Output |
|---|---|---|---|
| `caf validate` | CI | Process working directory | CAF report JSON and exit status |
| `caf validate --agent codex` | Codex `Stop` hook | Codex JSON on standard input | Codex `Stop` JSON on standard output |

The command accepts no project-root, validator-selection, OFT-path, or OFT-version option. Unknown commands, positional arguments, or agent values are usage failures.

## Project Discovery and Configuration

CAF searches for exactly one filename, `.caf.yaml`, starting at the invocation working directory and then visiting each parent directory. The nearest regular file wins. Reaching the filesystem root without finding the file is a tool failure. Supporting one filename and choosing the nearest match removes same-directory ambiguity.

The containing directory is the project root. Version 1 has this schema:

```yaml
version: 1
validators:
  - requirements-document
```

Unknown fields, unsupported versions, empty or duplicate validator names, and an empty validator list are invalid. Unknown validator identifiers produce a tool failure rather than being ignored.

## Shared Contracts

The normalized agent event contains `kind`, `working directory`, `session identifier`, and `already continued`. The initial kind is `turn_stop`. Provider fields do not enter validators.

Each validator returns one of:

- `pass`: the check completed and accepted the project state;
- `finding`: the check completed and found a rule violation;
- `tool_failure`: the check could not complete reliably; or
- `not_applicable`: policy determined that the check does not apply.

The aggregate order is `tool_failure`, `finding`, `pass`, then `not_applicable`. The agent decision is separate: it is either `allow_stop` or `continue_work` with a safe reason.

## Codex Stop Mapping

The adapter requires the documented Codex `Stop` fields `session_id`, `cwd`, `hook_event_name`, `turn_id`, and `stop_hook_active`. It ignores unknown fields for forward compatibility and requires `hook_event_name` to equal `Stop`.

| CAF decision | Codex JSON | Process exit |
|---|---|---:|
| `allow_stop` | `{}` | 0 |
| `continue_work` | `{"decision":"block","reason":"..."}` | 0 |

A finding, tool failure, invalid hook payload, or missing configuration requests continuation. When `stop_hook_active` is already true, CAF allows the turn to stop to prevent an uncontrolled loop; independent CI validation remains responsible for mandatory enforcement.

Hook standard output contains only protocol JSON. A serialization failure is written safely to standard error and exits 2 because CAF cannot provide a valid protocol response.

## CI Mapping

CI receives the complete CAF report as JSON. Exit status 0 means pass or not applicable, 1 means finding, and 2 means tool failure or invalid invocation. CI does not change validator semantics.

## Security and Compatibility

CAF limits hook input to one MiB, accepts exactly one JSON value, does not invoke a shell, and does not include prompts, credentials, transcript contents, or source contents in diagnostics. Configuration has an explicit version. Provider adapters are separate packages at the boundary and can evolve without changing shared validator results.

The Codex wire contract follows the [official Codex Hooks documentation](https://learn.chatgpt.com/docs/hooks). OFT execution remains dependent on issue #7; the validator router is designed to add it without changing `caf validate`.
