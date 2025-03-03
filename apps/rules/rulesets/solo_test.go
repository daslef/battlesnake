package rulesets

import (
	"rules"
	"rules/settings"
	"testing"

	"github.com/stretchr/testify/require"
)

func getSoloRuleset(settings settings.Settings) Ruleset {
	return NewRulesetBuilder().WithSettings(settings).NamedRuleset(rules.GameTypeSolo)
}

func TestSoloName(t *testing.T) {
	r := getSoloRuleset(settings.Settings{})
	require.Equal(t, "solo", r.Name())
}

func TestSoloCreateNextBoardStateSanity(t *testing.T) {
	boardState := &rules.BoardState{}
	r := getSoloRuleset(settings.Settings{})
	gameOver, _, err := r.Execute(boardState, []SnakeMove{})
	require.NoError(t, err)
	require.True(t, gameOver)
}

func TestSoloIsGameOver(t *testing.T) {
	tests := []struct {
		Snakes   []rules.Snake
		Expected bool
	}{
		{[]rules.Snake{}, true},
		{[]rules.Snake{{}}, false},
		{[]rules.Snake{{}, {}, {}}, false},
		{[]rules.Snake{{EliminatedCause: rules.EliminatedByOutOfBounds}}, true},
		{
			[]rules.Snake{
				{EliminatedCause: rules.EliminatedByOutOfBounds},
				{EliminatedCause: rules.EliminatedByOutOfBounds},
				{EliminatedCause: rules.EliminatedByOutOfBounds},
			},
			true,
		},
	}

	r := getSoloRuleset(settings.Settings{})
	for _, test := range tests {
		b := &rules.BoardState{
			Height: 11,
			Width:  11,
			Snakes: test.Snakes,
			Food:   []rules.Point{},
		}

		actual, _, err := r.Execute(b, nil)
		require.NoError(t, err)
		require.Equal(t, test.Expected, actual)
	}
}

// Checks that a single snake doesn't end the game
// also that:
// - snake moves okay
// - food gets consumed
// - snake grows and gets health from food
var soloCaseNotOver = gameTestCase{
	"Solo Case Game Not Over",
	&rules.BoardState{
		Width:  10,
		Height: 10,
		Snakes: []rules.Snake{
			{
				ID:     "one",
				Body:   []rules.Point{{X: 1, Y: 1}, {X: 1, Y: 2}},
				Health: 100,
			},
		},
		Food: []rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}},
	},
	[]SnakeMove{
		{ID: "one", Move: rules.MoveDown},
	},
	nil,
	&rules.BoardState{
		Width:  10,
		Height: 10,
		Snakes: []rules.Snake{
			{
				ID:     "one",
				Body:   []rules.Point{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 1}},
				Health: 100,
			},
		},
		Food: []rules.Point{{X: 0, Y: 0}},
	},
}

func TestSoloCreateNextBoardState(t *testing.T) {
	cases := []gameTestCase{
		// inherits these test cases from standard
		standardCaseErrNoMoveFound,
		standardCaseErrZeroLengthSnake,
		standardCaseMoveEatAndGrow,
		standardMoveAndCollideMAD,
		soloCaseNotOver,
	}
	r := getSoloRuleset(settings.Settings{})
	for _, gc := range cases {
		// test a RulesBuilder constructed instance
		gc.requireValidNextState(t, r)
		// also test a pipeline with the same settings
		gc.requireValidNextState(t, NewRulesetBuilder().PipelineRuleset(rules.GameTypeSolo, NewPipeline(soloRulesetStages...)))
	}
}

// Test a snake running right into the wall is properly eliminated
func TestSoloEliminationOutOfBounds(t *testing.T) {
	r := getSoloRuleset(settings.Settings{})

	// Using MaxRand is important because it ensures that the snakes are consistently placed in a way this test will work.
	// Actually random placement could result in the assumptions made by this test being incorrect.
	initialState, err := rules.CreateDefaultBoardState(rules.MaxRand, 2, 2, []string{"one"})
	require.NoError(t, err)

	_, next, err := r.Execute(initialState, []SnakeMove{{ID: "one", Move: "right"}})
	require.NoError(t, err)
	require.NotNil(t, initialState)

	ended, next, err := r.Execute(next, []SnakeMove{{ID: "one", Move: "right"}})
	require.NoError(t, err)
	require.NotNil(t, initialState)

	require.True(t, ended)
	require.Equal(t, rules.EliminatedByOutOfBounds, next.Snakes[0].EliminatedCause)
	require.Equal(t, "", next.Snakes[0].EliminatedBy)
	require.Equal(t, 1, next.Snakes[0].EliminatedOnTurn)
}
