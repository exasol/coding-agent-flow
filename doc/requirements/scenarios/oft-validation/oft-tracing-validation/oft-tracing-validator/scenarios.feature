Feature: OFT Tracing Validator
  Validate OFT tracing in a Target Project.

  @id:scn~oft-tracing-validator~1
  # Covers: req~oft-tracing-validation~1
  # Needs: at
  Scenario Outline: Validate OFT Requirement coverage
    Given a target project path containing the OFT Feature `<feature_id>` that needs an OFT Requirement
    And the target project contains `<coverage>`
    When the caller validates OFT tracing for the target project path
    Then the OFT Validation response has valid `<valid>`

    Examples:
      | feature_id                     | coverage                                                                            | valid |
      | feat~valid-fixture-feature~1   | the OFT Requirement `req~valid-fixture-requirement~1` that covers this OFT Feature | true  |
      | feat~invalid-fixture-feature~1 | no OFT Requirement that covers this OFT Feature                                     | false |
