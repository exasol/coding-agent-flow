---
name: reconcile-acceptance-tests
description: "Reconcile CAF acceptance tests with canonical Gherkin scenarios, use cases, and bounded-context domain models. Use only when the user says 'reconcile acceptance tests' or explicitly requests acceptance-test synchronization."
---

# Reconcile Acceptance Tests

Acceptance tests are executable evidence for the combined contract:

- Gherkin scenarios specify observable examples and acceptance outcomes.
- `use_cases.md` specifies the public actor, input, and outcome boundary.
- The bounded-context model, value objects, invariants, and failures define
  the legal vocabulary and rules.

None of these sources alone is sufficient. Do not derive acceptance tests from
Gherkin wording alone, and do not let test setup, mocks, private helpers, or
transport mechanics redefine the domain.

## Resolve the bounded context and test location

Treat each `.feature` file as a scenario list. Resolve its bounded context
from the first segment after `doc/requirements/scenarios/`:

```text
doc/requirements/scenarios/<bounded-context>/<feature>/<use-case>/scenarios.feature
doc/requirements/domain/<bounded-context>/
test/acceptance/<bounded-context>/<feature>/<go_use_case>_test.go
```

`<go_use_case>` is the `<use-case>` path name converted from hyphens to
underscores, because Go filenames use underscores. The current scenario list
therefore maps as follows:

| Artifact | Path |
| --- | --- |
| Gherkin | `doc/requirements/scenarios/spec-validation/spec-structure-validation/requirement-document-structure-validator/scenarios.feature` |
| Domain context | `doc/requirements/domain/spec-validation/` |
| Future acceptance test | `test/acceptance/spec-validation/spec-structure-validation/requirement_document_structure_validator_test.go` |

When several scenario lists are selected, group them by `<bounded-context>`
and reconcile each context independently. Do not infer a context from a
feature or use-case name when the first segment is absent or does not have a
matching domain directory; report that routing gap instead.

## Read the complete acceptance contract

For every selected context, read the global and context glossaries, context and
scope, entity model, value objects, invariants, use cases, failures, the
applicable system requirement, and every selected `.feature` file. Preserve
the requirement IDs, scenario IDs, and `Covers` links. Use trace-item
reconciliation before changing a named OFT identifier or revision.

Build this mapping before adding or editing tests:

| Scenario ID | Requirement | Use case | Public boundary | Test function |
| --- | --- | --- | --- | --- |
| `scn~...` | `req~...` from `Covers` | matching `<use-case>` directory | approved CLI or exported application use case | `Test<Outcome>` |

Generate the acceptance-test contract before the public boundary, production
code, or fixtures exist. For a new feature or bug fix, its acceptance evidence
is therefore normally expected to be red at first. Record the missing public
boundary and fixtures as the explicit reason for that red result; do not test
an internal package as a substitute for acceptance evidence.

An expected-red test must still compile and execute. Do not import a planned
Go package or refer directly to a planned type, function, method, or parameter
that does not yet exist: that only produces a compilation error, not useful
acceptance evidence. Instead, use the generic acceptance runner and its
scenario binding registry. Until a scenario has a public-boundary binding, the
runner fails that scenario with a clear missing-binding or missing-fixture
message. Do not use `t.Skip` for this condition.

## Acceptance framework for CAF

CAF retains OFT-tagged `.feature` files as the canonical BDD specification.
CAF's acceptance-test runner is a repository-owned Go test harness for CAF's
own public use cases. It resolves a canonical scenario ID, establishes the
scenario fixture, invokes the approved CAF public boundary once it exists, and
asserts the documented outcome. It is not a general-purpose, cross-language
execution framework.

Do not duplicate scenarios in JSON, Go structs, or a second feature file. The
runner must preserve their OFT IDs and invoke a CAF public CLI or exported
application use case rather than an internal package. Place black-box tests in
`test/acceptance/` and use Go's `testing` package and the existing `testify`
assertions.

The runner's binding registry maps each canonical scenario ID to its fixture,
action, and assertions, and fails the test when a binding does not yet exist.
A binding may name planned CAF public concepts in its metadata, but it may
only call Go symbols after they exist. This lets CAF acceptance tests be
generated and run before implementation without coupling them to an internal
package or making the whole suite fail to compile.

Add a repository-owned binding validation when acceptance tests are introduced.
It must verify that every selected `@id:scn~...` has exactly one corresponding
acceptance-test marker and that the marker points to that scenario. A planned
test marker has this shape:

```go
// [at~requirement-document-structure-is-valid~1 -> scn~requirement-document-structure-is-valid~1]
func TestRequirementDocumentStructureIsValid(t *testing.T) { /* ... */ }
```

The marker name identifies the acceptance-test evidence; the arrow binds it to
the canonical Gherkin scenario. The test then supplies the actual executable
evidence. Keep scenario text in `.feature` files; do not create a JSON or Go
struct shadow scenario model.

## Expected-red acceptance workflow

1. Document a feature or bug fix through its domain model, use case, and
   canonical Gherkin scenario.
2. Generate one OFT-tagged acceptance-test entry per scenario and register its
   planned public-boundary contract with the generic runner.
3. Run the suite. A missing binding, fixture, or public boundary fails the
   scenario with an explicit expected-red reason. This is the required
   pre-implementation state, not a readiness-only note.
4. Implement the public boundary, deterministic fixture, and binding. The
   same test remains red until the observable behavior satisfies its Then
   clauses.
5. Complete the work only when the acceptance test is green and its OFT
   traceability validation passes.

For a bug fix, first add or revise the scenario and expected-red acceptance
test that reproduces the defect; only then change the implementation. Do not
remove, skip, or weaken a failing acceptance assertion to make the suite pass.

## Translate scenario phases into tests

| Gherkin phase | Acceptance-test responsibility |
| --- | --- |
| Given | Once implemented, build only legal, deterministic fixture state through public inputs. For CAF file validation, create a temporary Source Project and its requirements documents. Before then, let the runner report the missing fixture as expected red. |
| When | Once implemented, invoke the approved public boundary exactly once for the scenario action. For a CLI, execute the built command; for an application boundary, call its exported use case. Before then, let the runner report the missing binding or boundary as expected red. |
| Then | Assert only the documented response, state, or failure outcome. Verify `RequirementStructure.valid` for the current Spec Validation scenarios. |

Use a composable fixture to establish a legal before-state, never to enforce an
invariant or replay a prior When/Then. Tests must be isolated: use a temporary
directory, controlled clock/inputs where relevant, no ordering dependence, no
shared mutable global state, and no sleep-based timing.

For valid and invalid scenario families, keep the same public invocation and
vary only the documented input. Assert the complete relevant outcome, not an
implementation detail such as parser call count, private type, cache state, or
filesystem traversal strategy.

## Reconcile tests with domain documents

| Source | Acceptance-test obligation |
| --- | --- |
| Context and glossary | Use the actor and terms exactly; do not invent test-only business synonyms. |
| Entity model and value objects | Build legal input values and observe public values, never mutate private state. |
| Invariants | Cover each valid/invalid branch through the public boundary; do not reimplement the rule in fixture helpers. |
| Use cases | Invoke the declared actor boundary and assert its documented outcome. |
| Failures | Assert only approved classifications and safe outputs; leave deferred diagnostics, serialization, and exit statuses untested until designed. |
| Gherkin | Bind one test to one scenario ID, preserve Given/When/Then order, and keep every observable assertion aligned with its Then clauses. |

If reconciliation changes the glossary, model, invariants, use case, or
scenario wording, invoke `$reconcile-domain-model`. If it changes a context
boundary, actor, request, response, or relationship shown in a diagram, invoke
`$reconcile-diagrams` too.

## Current readiness assessment

The current Spec Validation scenarios define `RequirementStructure.valid`, but
the contract explicitly defers the public CLI command, JSON schema,
diagnostics, and exit status. Generate the six scenario bindings now. Until a
public Requirement Document Structure Validator boundary, its deterministic
fixtures, and bindings exist, each affected acceptance test is expected red
with the runner's explicit missing-binding or missing-fixture failure. The
tests become green only after that public contract and the documented behavior
are implemented.

## Handoff

Report the scenario-to-context-to-test mapping, test coverage added or
missing, whether each test is intentionally red or green, the explicit red
reason, public-boundary readiness, command used to run the tests, and any
deferred contract gap. Run the repository's deterministic test and binding
validation when those commands exist.
