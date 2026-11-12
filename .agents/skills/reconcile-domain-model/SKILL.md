---
name: reconcile-domain-model
description: "Reconcile CAF Gherkin scenario lists with their bounded-context domain documentation. Use only when the user says 'reconcile domain model' or explicitly requests scenario-to-domain-model synchronization."
---

# Reconcile Domain Model

Reconcile Gherkin scenarios under `doc/requirements/scenarios/` with the
human-readable domain model under `doc/requirements/domain/`. The domain
documents are the source of domain meaning; Gherkin provides observable
behavioral evidence. Do not turn fixture details, parser mechanics, CLI
transport, or test implementation into domain concepts.

## Discover and partition scenario lists

Treat each `.feature` file as one scenario list. Resolve its bounded context
from the first path segment after `doc/requirements/scenarios/`:

```text
doc/requirements/scenarios/<bounded-context>/<feature>/<use-case>/scenarios.feature
                                  └───────────────┘
doc/requirements/domain/<bounded-context>/
```

For example, the current file
`doc/requirements/scenarios/spec-validation/spec-structure-validation/requirement-document-structure-validator/scenarios.feature`
belongs to the `spec-validation` bounded context and reconciles with
`doc/requirements/domain/spec-validation/`.

When given multiple scenario lists:

1. Inventory every `.feature` file and group it by `<bounded-context>`.
2. Reconcile one bounded context at a time, using all of that context's
   scenario lists before editing its documents.
3. Do not infer a bounded context from a feature or use-case directory. If the
   first segment is missing or no matching domain directory exists, report the
   path-to-context gap and request or use an explicitly authorized context
   name.
4. Update `doc/requirements/domain/glossary.md` only for vocabulary shared by
   two or more bounded contexts. Keep context-owned terms in the matching
   context glossary.

## Read the complete context before changing it

For each resolved `<bounded-context>`, read:

```text
doc/requirements/domain/index.md
doc/requirements/domain/glossary.md
doc/requirements/domain/<bounded-context>/index.md
doc/requirements/domain/<bounded-context>/glossary.md
doc/requirements/domain/<bounded-context>/context_and_scope.md
doc/requirements/domain/<bounded-context>/entity_model_relationship.md
doc/requirements/domain/<bounded-context>/value_objects.md
doc/requirements/domain/<bounded-context>/invariants.md
doc/requirements/domain/<bounded-context>/use_cases.md
doc/requirements/domain/<bounded-context>/failures.md
doc/requirements/scenarios/<bounded-context>/**/scenarios.feature
```

Read the applicable requirement in `doc/requirements/system_requirements.md`
and preserve its OFT identifier, scenario IDs, and `Covers` links. If a
semantic change requires renaming, superseding, or changing an OFT item
revision, stop and use the trace-item reconciliation workflow before changing
the identifier.

## Translate scenarios into the domain model

For every scenario, derive only actor-relevant facts:

| Scenario evidence | Domain translation |
| --- | --- |
| Actor issuing a request | Actor and use-case caller; define it in the context glossary. |
| Given state and named business nouns | Entity, value object, external source artifact, or legal precondition. Do not model test fixtures or file-parser internals. |
| When action and supplied values | One use case, its command/input boundary, and any value objects needed to express the input. |
| Then response or observable state | Response/read model attributes, state transition, or invariant outcome. |
| Multiple scenarios differing only in valid/invalid inputs | One use case and one invariant catalog with separate valid and failure branches. |

Apply these DDD distinctions consistently:

- An **entity** has a stable identity and lifecycle that matters to the actor.
- A **value object** is immutable and compared by value; it has no independent
  lifecycle. Put its validation boundary in `value_objects.md`.
- An **invariant** belongs to the entity or value object that can decide it;
  record its rule and the resulting response behavior in `invariants.md`.
- A **use case** coordinates the actor's request and returns the documented
  outcome; it does not own entity rules.
- A **bounded context** owns its vocabulary and model. Cross-context exchange
  uses an explicit identifier or response contract, never a live object graph.
- A **failure** belongs in `failures.md` only when the scenarios or approved
  requirements define its classification and outcome. Mark unapproved error,
  diagnostic, serialization, and exit-status behavior as deferred rather than
  inventing it.

## Synchronize each domain document

| Document | Reconcile with scenarios by ensuring |
| --- | --- |
| `glossary.md` | Every actor, entity, value object, field, response, and domain verb has one canonical term. Move only truly shared terms to the global glossary. |
| `context_and_scope.md` | Actor, external systems/artifacts, request boundary, response boundary, in-scope behavior, and deliberate exclusions match the scenario list. |
| `entity_model_relationship.md` | Entities, identity, cardinality, ownership, optionality, and response model explain every scenario noun and observable result. |
| `value_objects.md` | Immutable inputs and fields have values, owners, and validation boundaries consistent with Given/When values. |
| `invariants.md` | Every valid/invalid branch has a named rule, owner, and result calculation. Do not restate the full scenario text. |
| `use_cases.md` | Each `<use-case>` scenario-path segment maps to one named use case, describes its actor/input/outcome, and links to its canonical `.feature` file. |
| `failures.md` | Only approved scenario/requirement failures are specified; unsupported filesystem, diagnostics, transport, and exit behavior remains deferred. |
| `index.md` | Lists every context document that now carries domain meaning. |

If the entity model, context boundary, actor, input, outcome, relationship, or
invariant changes, invoke `$reconcile-diagrams` and synchronize the context,
model/ER, and use-case diagrams with the updated documents.

## Reconciliation sequence and handoff

1. Produce a mapping of scenario path, requirement ID, bounded context, and
   use-case name.
2. Identify terminology, entity/value-object, invariant, use-case, and
   failure gaps before editing.
3. Update glossary and model documents before updating diagrams.
4. Keep each scenario list in its resolved context; do not duplicate scenarios
   in Markdown.
5. Verify each scenario's actor, input, outcome, and valid/invalid branch is
   represented by the context documents.
6. Report changed paths, the scenario-to-context mapping, and any deferred or
   unresolved design gap.
