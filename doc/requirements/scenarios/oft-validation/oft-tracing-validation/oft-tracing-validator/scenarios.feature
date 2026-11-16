Feature: OFT Tracing Validator
  Validate OFT tracing in a Target Project.

  @id:scn~oft-requirement-covers-feature~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario: Return valid when an OFT Requirement covers the OFT Feature that needs it
    Given a target project path containing `features.md` with the OFT Feature `feat~valid-fixture-feature~1` that needs an OFT Requirement
    And `requirements.md` with the OFT Requirement `req~valid-fixture-requirement~1` that covers `feat~valid-fixture-feature~1`
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response has valid `true`

  @id:scn~oft-requirement-does-not-cover-feature~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario: Return invalid when an OFT Feature that needs an OFT Requirement stays uncovered
    Given a target project path containing `features.md` with the OFT Feature `feat~invalid-fixture-feature~1` that needs an OFT Requirement
    And `feat~invalid-fixture-feature~1` stays uncovered
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response has valid `false`
