Feature: Turn Decision
  The Turn Decision ends or continues an Agent Turn from the validation and the Continuation Active flag.

  @id:scn~domain-turn-decision_end-turn-when-validation-passes~1
  # Covers: dom~turn-decision~1
  # Needs: at
  Scenario Outline: End the turn when validation passes
    Given a Validation Report with the Overall Result `<overall result>`
    And Continuation Active is not set
    When the turn is decided
    Then the turn ends without reason

    Examples:
      | overall result |
      | pass           |
      | not_applicable |

  @id:scn~domain-turn-decision_continue-turn-on-problems~1
  # Covers: dom~turn-decision~1
  # Needs: at
  Scenario: Continue the turn with the reason when validation finds a problem
    Given a Validation Report with the results `[a: finding, b: pass, c: tool_failure]`
    And Continuation Active is not set
    When the turn is decided
    Then the turn continues
    And the reason names the results `a` and `c` and not `b`

  @id:scn~domain-turn-decision_continue-turn-on-validation-error~1
  # Covers: dom~turn-decision~1
  # Needs: at
  Scenario: Continue the turn with the error code when validation fails
    Given a Validation Runner error with the code `E-CAF-5`
    And Continuation Active is not set
    When the turn is decided
    Then the turn continues
    And the reason contains `E-CAF-5`

  @id:scn~domain-turn-decision_end-turn-when-continuation-active~1
  # Covers: dom~turn-decision~1
  # Needs: at
  Scenario Outline: End the turn and keep the reason when Continuation Active is set
    Given <validation>
    And Continuation Active is set
    When the turn is decided
    Then the turn ends
    And the reason contains `<reason>`

    Examples:
      | validation                                          | reason     |
      | a Validation Report with the results `[a: pass]`    |            |
      | a Validation Report with the results `[a: finding]` | a: finding |
      | a Validation Runner error with the code `E-CAF-5`   | E-CAF-5    |
