Feature: Agent Turn Validator
  Validate the Target Project when a coding agent is about to end its turn.

  @id:scn~agent-turn-validator_end-turn~1
  # Covers: req~agent-turn-validation~1
  # Needs: at
  Scenario: End the turn when validation passes
    Given a target project path whose `.caf.yaml` selects `oft-tracing-validator`
    And the OFT tracing of the target project is `valid`
    When the agent stops with a Stop Event without Continuation Active
    Then the TurnDecision response has action `end_turn`

  @id:scn~agent-turn-validator_continue-turn~1
  # Covers: req~agent-turn-validation~1
  # Needs: at
  Scenario Outline: Continue the turn when validation finds a problem or fails
    Given a target project path whose `.caf.yaml` selects `<validation>`
    And the OFT tracing of the target project is `<tracing>`
    When the agent stops with a Stop Event without Continuation Active
    Then the TurnDecision response has action `continue_turn`
    And the reason contains `<reason>`

    Examples:
      | validation            | tracing | reason                  |
      | oft-tracing-validator | invalid | CAF-OFT-TRACING-INVALID |
      | unknown-check         | valid   | CAF-UNKNOWN-VALIDATOR   |
