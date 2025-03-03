package rulesets

import (
	"rules"
	"testing"

	"github.com/stretchr/testify/require"
)

type gameTestCase struct {
	name          string
	prevState     *rules.BoardState
	moves         []SnakeMove
	expectedError error
	expectedState *rules.BoardState
}

// requireValidNextState requires that the ruleset produces a valid next state
func (gc *gameTestCase) requireValidNextState(t *testing.T, r Ruleset) {
	t.Helper()
	t.Run(gc.name, func(t *testing.T) {
		t.Helper()
		prev := gc.prevState.Clone() // clone to protect against mutation (so we can re-use test cases)
		_, nextState, err := r.Execute(prev, gc.moves)
		require.Equal(t, gc.expectedError, err)
		if gc.expectedState != nil {
			require.Equal(t, gc.expectedState.Width, nextState.Width)
			require.Equal(t, gc.expectedState.Height, nextState.Height)
			require.Equal(t, gc.expectedState.Food, nextState.Food)
			require.Equal(t, gc.expectedState.Snakes, nextState.Snakes)
		}
	})
}

func mockSnakeMoves() []SnakeMove {
	return []SnakeMove{
		{ID: "test-mock-move", Move: "mocked"},
	}
}
