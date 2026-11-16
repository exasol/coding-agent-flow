---
name: reconcile-code-implementation
description: Reconcile changed acceptance contracts and implement them in one agent workflow. Use when a task requires the red acceptance phase followed by the green production phase.
---

# Reconcile Code Implementation

Run acceptance reconciliation and implementation in one agent workflow.

## Change scope by pull-request type

Use the listed artifact categories as the pull-request change scope.

- Feature: requirements, glossary, use-case, and scenario documents; acceptance tests (`at_data.toml`, `at_test.go`); public API, application, integration, and unit-test code; trace markers.
- Bug fix: regression scenario, its `at_data.toml` blocks and checks; public API, application, integration, and unit-test code required by the regression; trace markers.
- Refactoring: implementation and unit tests that preserve approved behavior.
- Release: release notes, changelog, version sources, release configuration, and release validation.
- Security: security requirements and scenarios, security acceptance tests, affected source or dependency configuration, security documentation, and trace markers.
- Documentation: affected Markdown and cross-reference files.

## Feature pull requests

1. Apply `reconcile-acceptance-tests`.
2. Continue with `reconcile-acceptance-implementation` after acceptance reconciliation completes.
3. Write the production code.
4. Run acceptance tests until they pass.
5. Run configured traceability validation.
6. Report the completed contract, implementation, and validation results.

## Bug-fix pull requests

1. Apply `reconcile-acceptance-tests` to create the red regression contract.
2. Continue with `reconcile-acceptance-implementation` after reconciliation completes.
3. Write the production code for the preserved regression contract.
4. Run the regression acceptance test until it passes.
5. Run configured traceability validation.
6. Report the regression result and validation results.

## Refactoring pull requests

1. Run existing acceptance tests and verify a green result.
2. Refactor the production code.
3. Run acceptance tests and unit tests.
4. Run configured traceability validation.
5. Report the preserved behavior and validation results.

## Release pull requests

1. Update release notes, changelog, version sources, and release configuration.
2. Run release validation and configured traceability validation.
3. Report the release version and validation results.

## Security pull requests

1. Run existing acceptance tests and verify a green result.
2. Update the affected source or dependency configuration and security documentation.
3. Run security, acceptance, unit, and configured traceability validation.
4. Report the security behavior and validation results.

## Documentation pull requests

1. Update the affected documentation and cross-reference files.
2. Run documentation validation and configured traceability validation.
3. Report the updated documentation and validation results.

## Clarification and validation

Return acceptance-reconciliation questions, implementation questions, and validation results to the user. Continue with the next phase after the current phase completes.
