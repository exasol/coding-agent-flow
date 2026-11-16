---
name: reconcile-acceptance-implementation
description: Reconcile approved changed scenarios, acceptance tests, and public API with production implementation. Use for the green implementation phase after acceptance evidence provides a red result.
---

# Reconcile Acceptance Implementation

Implement the acceptance contract after `reconcile-acceptance-tests` completes. Use `reconcile-code-implementation` for the combined workflow.

## Feature pull requests

1. Preserve the red acceptance contract.
2. Implement the required public API, application behavior, integration behavior, and response mapping.
3. Use the shared public response type for related criteria that call the same public method.
4. Add focused unit tests for behavior behind the public API.
5. Run affected acceptance tests until they pass.
6. Run unit tests and repository validation commands.
7. Run configured traceability validation.
8. Report public behavior, production changes, and validation results.

## Bug-fix pull requests

1. Start with the red regression contract.
2. Preserve its scenario, `at_data.toml` blocks with their `id` and `description`, and checks.
3. Implement the behavior required by the regression scenario.
4. Run the regression acceptance test until it passes.
5. Run related unit tests and configured traceability validation.
6. Report the corrected observable behavior and validation results.

## Refactoring pull requests

1. Keep approved scenarios, `at_data.toml` blocks with their `id` and `description`, and checks as the behavior contract.
2. Refactor public, application, and integration code behind that contract.
3. Run affected acceptance tests until they pass.
4. Run unit tests, repository validation commands, and configured traceability validation.
5. Report the preserved behavior and validation results.

## Clarification

Pause implementation and ask Socratic questions when scenario behavior, fixture meaning, public method signatures, response fields, or error results need a product decision. Resume with the approved decision.
