---
name: reconcile-acceptance-tests
description: Reconcile changed Gherkin scenarios with the acceptance tests and the planned public API. Use for the red acceptance-contract phase before production implementation.
---

# Reconcile Acceptance Tests

Create a complete, executable acceptance contract before production implementation.

## Read the complete acceptance contract

For each selected feature, read:

- `<requirements-root>/glossary.md`
- `<requirements-root>/<feature>/glossary.md`
- `<requirements-root>/<feature>/use_cases.md`
- `<scenario-root>/<feature>/<requirement>/<use-case>/scenarios.feature`
- `test/acceptance/<feature>/<requirement>/<use-case>/at_test.go` and `at_data.toml`
- the shared When helper in `test/acceptance/<feature>/`
- `<public-api>.go`
- `internal/<domain>/` (domain, application port and service, adapter)

## Feature pull requests

1. Document the feature in its domain context, requirement, use case, and canonical Gherkin scenarios.
2. Write the acceptance test (see "Acceptance Test Files").
3. Run the affected acceptance suite and record its red result.
4. Hand the complete contract to `reconcile-acceptance-implementation`.

## Bug-fix pull requests

1. Add or revise the canonical Gherkin scenario that reproduces the defect.
2. Write the regression acceptance test (see "Acceptance Test Files").
3. Run it and record its red result.
4. Keep its checks as the regression contract.
5. Hand the regression contract to `reconcile-acceptance-implementation`.

## Refactoring pull requests

1. Run the affected acceptance suite before the refactoring.
2. Keep each scenario, `at_data.toml` block, and check as the executable behavior contract.
3. Run the affected acceptance suite and record its result.

## Acceptance Test Files

The Gherkin scenarios document the use case; the acceptance test carries its data. Write the test from the scenarios, as "Acceptance Tests" in `doc/developer_guide.md` shows:

- Scenario text: wrap each parameter value of a Given, When, or Then step in backticks, for example `` `items/item.txt` `` or `` `<item_id>` ``.
- `at_data.toml`: one `[[scenario]]` block per Scenario, and one per Examples row of a Scenario Outline. Each block starts with:
  - `id`: the scenario `@id` without the `@id:` prefix, for example `id = "scn~list-items~1"`. All blocks of one Scenario Outline share its `id`.
  - `description`: the scenario title after `Scenario:` or `Scenario Outline:`. For an Examples row, append `: ` and the cell that tells the rows apart, for example `description = "List the items of a shop: main"`.
- Put the values and files of the Given lines into `[scenario.given]` and `[[scenario.given.file]]`, and the expected values of the Then lines into `[scenario.then]`. For a value already in `[scenario.given]`, write `{{name}}` in `[scenario.then]`.
- `at_test.go`: the When line, the public method it calls, one check per kind of Then line, and one `[at -> scn~...]` marker per scenario `@id`. A check receives the response, Given values, and Then values.
- Keep the scenario and the test in sync: every file, ID, and value in `at_data.toml` appears in the scenario's Given and Then lines or Examples rows, every concrete input or expected value of the scenario appears in `at_data.toml`, the When line in `at_test.go` equals the scenario's When line, and each block's `id` and `description` follow the scenario's current `@id` and title.
- Run the configured traceability validation to confirm each scenario `@id` has its `at` marker.

## Record the Red Result

The fixture in `at_test.go` stores the public call as text, so the acceptance-test package compiles without compiling that call. During the test, the runner builds a small program for the call. Fixture compile errors report public-API gaps; `Then:` errors report behavior differences.

## Clarification

Pause reconciliation and ask Socratic questions when the contract needs a product or API decision. Ask for the governing scenario, fixture state, public method, response shape, or observable error result. Resume with the approved decision.
