---
name: reconcile-diagrams
description: "Reconcile human-readable Mermaid diagrams with CAF domain documentation. Use only when the user says 'reconcile diagrams' or explicitly requests a diagram change."
---

# Reconcile Diagrams

Reconcile the smallest diagrams that allow a named reader to answer a specific
question about the affected CAF domain. Diagrams are communicative evidence;
they do not replace requirements, domain text, Gherkin scenarios, or tests.

## Scope and location

Read `doc/requirements/domain/index.md`, the affected bounded-context index,
its glossary, model, use cases, invariants, and acceptance scenarios before
changing a diagram. Keep each diagram in its owning context:

| Reader question | Canonical location |
| --- | --- |
| Who is outside the bounded context and what crosses its boundary? | `context_and_scope.md` |
| What entities, value objects, and relationships make up the model? | `entity_model_relationship.md` |
| How does an actor invoke a use case and receive its outcome? | `use_cases.md` |

Use Mermaid fenced blocks in the owning Markdown document. Do not create image
assets, duplicate a diagram in several documents, add a class diagram merely
because code exists, or introduce implementation, CLI, persistence, or event
semantics that the documented domain does not establish.

## Recovery workflow

1. State the reader question and the intended outcome in adjacent prose.
2. Reconcile the diagram with the bounded-context glossary, entity model,
   invariants, use case, and Gherkin scenarios. Use exact glossary terms.
3. Use one abstraction level per diagram. Keep the main reading flow
   left-to-right or top-to-bottom; split a crowded view instead of accepting
   crossings or verbose labels.
4. Label directed relationships with concise verb phrases. Show a request and
   its response separately when both matter. State cardinality and ownership
   in an ER diagram or its adjacent prose.
5. Add a brief text alternative immediately after each diagram. Use colour or
   icons only as reinforcement, never as the sole source of meaning.
6. Render with the repository's configured Mermaid command when one exists;
   do not download or add an unpinned renderer. Otherwise inspect the rendered
   Markdown in the normal documentation viewer and report that no deterministic
   renderer is configured.

## Validation

Confirm all Mermaid blocks are syntactically present, links still resolve, and
the changed diagrams agree with their prose and Gherkin scenarios. Report the
rendering result separately from this structural review.

## Synchronize diagrams with domain documents

Treat the domain documents as the source of meaning and diagrams as focused
views of that meaning. Before changing a diagram, read the corresponding rows
below; after changing one of those documents, review the mapped diagram for
drift. Do not make a diagram the sole source of a rule or relationship.

| Diagram | Source documents to synchronize | Required agreement |
| --- | --- | --- |
| Context diagram in `context_and_scope.md` | `glossary.md`, `context_and_scope.md`, `use_cases.md` | Actor, external Source Project, use-case boundary, request, and response use the same terms and scope. |
| Domain model and ER diagrams in `entity_model_relationship.md` | `glossary.md`, `entity_model_relationship.md`, `value_objects.md`, `invariants.md` | Every shown entity and value object exists in the model; ownership, cardinality, optionality, and invariant owner match the prose. |
| Use-case flow in `use_cases.md` | `glossary.md`, `use_cases.md`, `invariants.md`, the linked `.feature` file | Actor, input, outcome, and decision rules match the use case and every relevant Gherkin scenario. |

Apply these checks whenever the corresponding source changes:

1. A glossary rename updates every diagram label and text alternative in that
   bounded context.
2. An entity, value object, relationship, cardinality, or invariant-owner
   change updates the domain model/ER diagram and its adjacent alternative.
3. An actor, input, outcome, or scenario change updates the use-case flow;
   update the context diagram too if a boundary crossing changes.
4. A context scope change updates the context diagram only when its actor,
   external dependency, request, or response changes.

Finish by comparing diagram labels and edges against the current documents,
then report the synchronized paths and any deliberately omitted view.
