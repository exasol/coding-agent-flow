Feature: Requirement Document Structure Validator
  Validate the structural invariants of Requirement Documents in a supplied
  Source Project.

  @id:scn~requirement-structure-is-valid~1
  # Covers: req~requirement-structure-validation~1
  Scenario: Return valid when every OFT requirement has Needs and Rationale
    Given a source project path containing a requirements document
    And the document contains OFT requirement items with non-empty Needs and Rationale fields
    When the caller validates the spec structure for the source project path
    Then the RequirementStructure response has valid true

  @id:scn~requirement-structure-missing-needs~1
  # Covers: req~requirement-structure-validation~1
  Scenario: Return invalid when an OFT requirement omits Needs
    Given a source project path containing a requirements document
    And an OFT requirement item without a Needs field
    When the caller validates the spec structure for the source project path
    Then the RequirementStructure response has valid false

  @id:scn~requirement-structure-blank-needs~1
  # Covers: req~requirement-structure-validation~1
  Scenario: Return invalid when an OFT requirement has blank Needs
    Given a source project path containing a requirements document
    And an OFT requirement item whose Needs field is empty or whitespace only
    When the caller validates the spec structure for the source project path
    Then the RequirementStructure response has valid false

  @id:scn~requirement-structure-missing-rationale~1
  # Covers: req~requirement-structure-validation~1
  Scenario: Return invalid when an OFT requirement omits Rationale
    Given a source project path containing a requirements document
    And an OFT requirement item without a Rationale field
    When the caller validates the spec structure for the source project path
    Then the RequirementStructure response has valid false

  @id:scn~requirement-structure-blank-rationale~1
  # Covers: req~requirement-structure-validation~1
  Scenario: Return invalid when an OFT requirement has blank Rationale
    Given a source project path containing a requirements document
    And an OFT requirement item whose Rationale field is empty or whitespace only
    When the caller validates the spec structure for the source project path
    Then the RequirementStructure response has valid false

  @id:scn~requirement-structure-no-requirements~1
  # Covers: req~requirement-structure-validation~1
  Scenario: Return valid when no OFT requirement items are present
    Given a source project path containing requirements documents without OFT requirement items
    When the caller validates the spec structure for the source project path
    Then the RequirementStructure response has valid true
