Feature: OFT Tracing Validator
  Validate OFT tracing in a Target Project.

  @id:scn~oft-requirement-covers-feature~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario: Return valid when every OFT Requirement covers an OFT Feature
    Given a target project path containing a requirements document
    And the document contains an OFT Feature item
    And every OFT Requirement item covers the OFT Feature item
    When the caller validates OFT tracing for the target project path
    Then the OFTValidation response has valid true

  @id:scn~oft-requirement-does-not-cover-feature~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario: Return false when an OFT Requirement awaits an OFT Feature link
    Given a target project path containing a requirements document
    And the document contains an OFT Requirement item awaiting an OFT Feature link
    When the caller validates OFT tracing for the target project path
    Then the OFTValidation response has valid false
