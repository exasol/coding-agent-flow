---
name: reconcile-code-implementation
description: Reconcile changed acceptance contracts and implement them in one agent workflow. Use when a task requires the red acceptance phase followed by the green production phase.
---

# Reconcile Code Implementation

Run acceptance reconciliation and implementation in one agent workflow.

## Change scope by pull-request type

Use the listed artifact categories as the pull-request change scope.

- Feature: requirements, glossary, use-case, and scenario documents; acceptance fixtures and bindings; public API, domain, adapter, and unit-test code; OFT trace markers.
- Bug fix: regression scenario, fixture, binding, and assertion; public API, domain, adapter, and unit-test code required by the regression; OFT trace markers.
- Refactoring: implementation, unit tests, acceptance bindings that adapt an approved public API, and related code mappings.
- Release: `doc/changes/`, changelog, version sources, release configuration, and release validation.
- Security: security requirements and scenarios, security fixtures and bindings, affected source or dependency configuration, security documentation, and trace markers.
- Documentation: affected Markdown, reStructuredText, glossary, use-case, and cross-reference files.

## Feature pull requests

1. Apply `reconcile-acceptance-tests`.
2. Continue with `reconcile-acceptance-implementation` after acceptance reconciliation reports success.
3. Write the production code.
4. Run acceptance tests until they pass.
5. Verify requirement traceability with OFT.
6. Report the completed contract, implementation, and validation results.

## Bug-fix pull requests

1. Apply `reconcile-acceptance-tests` to create the red regression contract.
2. Continue with `reconcile-acceptance-implementation` after reconciliation reports success.
3. Write the production code for the preserved regression contract.
4. Run the regression acceptance test until it passes.
5. Verify requirement traceability with OFT.
6. Report the regression result and validation results.

## Refactoring pull requests

1. Apply `reconcile-acceptance-tests` to establish the current behavior contract.
2. Continue with `reconcile-acceptance-implementation` after reconciliation reports success.
3. Write the refactored production code.
4. Run acceptance tests until they pass.
5. Verify requirement traceability with OFT.
6. Report the preserved behavior and validation results.

## Release pull requests

1. Update release notes, changelog, version sources, and release configuration.
2. Run release validation and OFT traceability verification.
3. Report the release version and validation results.

## Security pull requests

1. Apply `reconcile-acceptance-tests` to security behavior that has a public boundary.
2. Update the affected source or dependency configuration and security documentation.
3. Run security, acceptance, unit, and OFT traceability validation.
4. Report the security behavior and validation results.

## Documentation pull requests

1. Update the affected documentation, glossary terms, use cases, and cross-references.
2. Run documentation validation and OFT traceability verification.
3. Report the updated documentation and validation results.

## Clarification and validation

Return acceptance-reconciliation questions, implementation questions, and validation results to the user. Continue with the next phase after the current phase reports success.
