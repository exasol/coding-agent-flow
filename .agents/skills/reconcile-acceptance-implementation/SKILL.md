---
name: reconcile-acceptance-implementation
description: Reconcile approved changed acceptance scenarios, fixtures, bindings, and public API with production implementation. Use for the green implementation phase after acceptance evidence provides a red result.
---

# Reconcile Acceptance Implementation

Implement the acceptance contract after `reconcile-acceptance-tests` reports success. Use `reconcile-code-implementation` for the combined workflow.

## Feature pull requests

1. Preserve the red acceptance contract.
2. Implement the required public endpoint, domain behavior, adapter behavior, and response mapping.
3. Use the shared public response type for related criteria that call the same endpoint.
4. Add focused unit evidence for domain behavior behind the public boundary.
5. Run affected acceptance tests until they pass.
6. Run unit tests and repository validation commands.
7. Verify requirement traceability with OFT.
8. Report public behavior, production changes, and validation results.

## Bug-fix pull requests

1. Start with the red regression contract.
2. Preserve its scenario, fixture, binding, and assertion.
3. Implement the behavior required by the regression scenario.
4. Run the regression acceptance test until it passes.
5. Run related unit tests and verify requirement traceability with OFT.
6. Report the corrected observable behavior and validation results.

## Refactoring pull requests

1. Keep approved scenarios, fixtures, bindings, and expected responses as the behavior contract.
2. Refactor public, domain, and adapter code behind that contract.
3. Run affected acceptance tests until they pass.
4. Run unit tests, repository validation commands, and OFT traceability verification.
5. Report the preserved behavior and validation results.

## Clarification

Pause implementation and ask Socratic questions when scenario behavior, fixture meaning, public method signatures, response fields, or error results need a product decision. Resume with the approved decision.
