---
name: reconcile-acceptance-tests
description: Reconcile changed Gherkin scenarios with acceptance bindings, fixtures, and planned public API. Use for the red acceptance-contract phase before production implementation.
---

# Reconcile Acceptance Tests

Create a complete, executable acceptance contract before production implementation.

## Read the complete acceptance contract

For each selected context, read:

- `doc/requirements/domain/glossary.md`
- `doc/requirements/domain/<context>/glossary.md`
- `doc/requirements/domain/<context>/use_cases.md`
- `doc/requirements/features/features.md`
- `doc/requirements/features/requirements.md`
- `doc/requirements/scenarios/<feature>/<requirement>/<use-case>/scenarios.feature`
- `test/acceptance/<use-case>_bindings_test.go`
- `test/acceptance/fixtures/<feature>/<requirement>/<use-case>.fixture`
- `oft_validation.go`
- `internal/oftvalidation/*.go`
- `internal/openfasttrace/*.go`

## Feature pull requests

1. Document the feature in the domain model, use case, and canonical Gherkin scenario.
2. Add one `at` marker, public-boundary program, deterministic fixture, and typed binding for each scenario.
3. Write the acceptance test.
4. Run the affected acceptance suite and record its red result.
5. Hand the complete contract to `reconcile-acceptance-implementation`.

## Bug-fix pull requests

1. Add or revise the canonical Gherkin scenario that reproduces the defect.
2. Add one `at` marker, public-boundary program, deterministic fixture, and typed binding for the regression.
3. Write the regression acceptance test.
4. Run it and record its red result.
5. Preserve the acceptance assertion as the regression contract.
6. Hand the regression contract to `reconcile-acceptance-implementation`.

## Refactoring pull requests

1. Run the affected acceptance suite before the refactoring.
2. Keep each scenario, fixture, binding, and expected response as the executable behavior contract.
3. Update the public-boundary binding when the approved public API changes.
4. Run the affected acceptance suite and record its result.

## Fixtures and bindings

- Gherkin scenario files own canonical scenario IDs.
- Fixture catalog files own Target Project files and content.
- Go binding functions register each canonical scenario ID with typed expected response values.
- The Go runner owns the runtime binding registry and resolves scenario IDs, fixtures, and bindings.
- The Go runner reports the reason for each pending binding, fixture, or public boundary.
- Store Target Project state in external fixture files.
- Use one blank line between consecutive fixture files and scenarios.
- Map `Given` to fixture state, `When` to one public call, and `Then` to typed response values.
- Compare OFT fixture element IDs with the matching response element-type list.
- Express each unrecognized element type with an empty response list.

## Clarification

Pause reconciliation and ask Socratic questions when the contract needs a product or API decision. Ask for the governing scenario, fixture state, public endpoint, response shape, element IDs, or observable error result. Resume with the approved decision.
