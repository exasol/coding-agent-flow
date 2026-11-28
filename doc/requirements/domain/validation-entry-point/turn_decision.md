# Turn Decision

`dom~turn-decision~1`

The Turn Decision tells a coding agent whether it may end its Agent Turn. It decides from the Validation Report, or from a Validation Runner error, and the Continuation Active flag of the Stop Event, and holds the action `end_turn` or `continue_turn`, the reason, and the Validation Report.

Needs: scn, impl
