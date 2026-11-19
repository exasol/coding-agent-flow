Feature: Validation Runner
  Run the validations that a Target Project selects.

  @id:scn~validation-runner_selected-validations~1
  # Covers: req~validation-selection~1
  # Needs: at
  Scenario Outline: Run each selected validation
    Given a target project path whose `.caf.yaml` selects `[oft-tracing-validator]`
    And the OFT tracing of the target project is `<tracing>`
    When the caller runs the validations of the target project path
    Then the ValidationReport response lists the results `<results>`

    Examples:
      | tracing | results                          |
      | valid   | [oft-tracing-validator: pass]    |
      | invalid | [oft-tracing-validator: finding] |

  @id:scn~validation-runner_unknown-validation~1
  # Covers: req~validation-selection~1
  # Needs: at
  Scenario: Report a tool failure for an unknown validation
    Given a target project path whose `.caf.yaml` selects `[oft-tracing-validator, unknown-check]`
    And the OFT tracing of the target project is `valid`
    When the caller runs the validations of the target project path
    Then the ValidationReport response lists the results `[oft-tracing-validator: pass, unknown-check: tool_failure]`
    And the result of `unknown-check` has code `CAF-UNKNOWN-VALIDATOR`

  @id:scn~validation-runner_overall-result~1
  # Covers: req~validation-result~1
  # Needs: at
  Scenario Outline: Report the most severe result as overall result
    Given a target project path whose `.caf.yaml` selects `<validations>`
    And the OFT tracing of the target project is `<tracing>`
    When the caller runs the validations of the target project path
    Then the ValidationReport response has status `<overall result>`

    Examples:
      | validations                            | tracing | overall result |
      | [oft-tracing-validator]                | valid   | pass           |
      | [oft-tracing-validator]                | invalid | finding        |
      | [oft-tracing-validator, unknown-check] | invalid | tool_failure   |
