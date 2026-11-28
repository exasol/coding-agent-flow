Feature: Agent Continuation Limiter
  Let a coding agent end its turn when the turn already continued because of CAF.

  @id:scn~agent-continuation-limiter~1
  # Covers: req~agent-continuation-limit~1
  # Needs: at
  Scenario: End the turn when the turn already continued because of CAF
    Given a target project path whose `.caf.yaml` selects `oft-tracing-validator`
    And the OFT tracing of the target project is `invalid`
    When the agent stops with a Stop Event with Continuation Active
    Then the TurnDecision response has action `end_turn`
    And the reason contains `CAF-OFT-TRACING-INVALID`
