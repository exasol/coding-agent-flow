Feature: OFT Tracing Validator
  Validate OFT tracing in a Target Project.

  @id:scn~oft-tracing-validator_valid~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario: Return valid when an OFT Requirement covers the OFT Feature that needs it
    Given a target project path containing `features.md` with the OFT Feature `feat~valid-fixture-feature~1` that needs an OFT Requirement
    And `requirements.md` with the OFT Requirement `req~valid-fixture-requirement~1` that covers `feat~valid-fixture-feature~1`
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response has valid `true`

  @id:scn~oft-tracing-validator_invalid~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario: Return invalid when an OFT Feature that needs an OFT Requirement stays uncovered
    Given a target project path containing `features.md` with the OFT Feature `feat~invalid-fixture-feature~1` that needs an OFT Requirement
    And `feat~invalid-fixture-feature~1` stays uncovered
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response has valid `false`

  @id:scn~oft-tracing-validator_requirement-covers-other-feature~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario: Return invalid when the OFT Requirement covers another OFT Feature than the one that needs it
    Given a target project path containing `features.md` with the OFT Feature `feat~needed-fixture-feature~1` that needs an OFT Requirement
    And `requirements.md` with the OFT Requirement `req~other-feature-fixture-requirement~1` that covers `feat~other-fixture-feature~1`
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response has valid `false`
