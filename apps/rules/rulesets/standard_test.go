package rulesets

import (
	"math"
	"math/rand"
	"rules"
	"rules/settings"
	"testing"

	"github.com/stretchr/testify/require"
)

func getStandardRuleset(settings settings.Settings) Ruleset {
	return NewRulesetBuilder().WithSettings(settings).NamedRuleset(rules.GameTypeStandard)
}

func TestSanity(t *testing.T) {
	r := getStandardRuleset(settings.Settings{})

	state, err := rules.CreateDefaultBoardState(rules.MaxRand, 0, 0, []string{})
	require.NoError(t, err)
	require.NotNil(t, state)

	gameOver, state, err := r.Execute(state, []SnakeMove{})
	require.NoError(t, err)
	require.True(t, gameOver)
	require.NotNil(t, state)
	require.Equal(t, 0, state.Width)
	require.Equal(t, 0, state.Height)
	require.Len(t, state.Food, 0)
	require.Len(t, state.Snakes, 0)
}

func TestStandardName(t *testing.T) {
	r := getStandardRuleset(settings.Settings{})
	require.Equal(t, "standard", r.Name())
}

// Checks that the error for a snake missing a move is returned
var standardCaseErrNoMoveFound = gameTestCase{
	"Standard Case Error No Move Found",
	&rules.BoardState{
		Width:  10,
		Height: 10,
		Snakes: []rules.Snake{
			{
				ID:     "one",
				Body:   []rules.Point{{X: 1, Y: 1}, {X: 1, Y: 2}},
				Health: 100,
			},
			{
				ID:     "two",
				Body:   []rules.Point{{X: 3, Y: 4}, {X: 3, Y: 3}},
				Health: 100,
			},
		},
		Food: []rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}},
	},
	[]SnakeMove{
		{ID: "one", Move: rules.MoveUp},
	},
	rules.ErrorNoMoveFound,
	nil,
}

// Checks that the error for a snake with no points is returned
var standardCaseErrZeroLengthSnake = gameTestCase{
	"Standard Case Error Zero Length Snake",
	&rules.BoardState{
		Width:  10,
		Height: 10,
		Snakes: []rules.Snake{
			{
				ID:     "one",
				Body:   []rules.Point{{X: 1, Y: 1}, {X: 1, Y: 2}},
				Health: 100,
			},
			{
				ID:     "two",
				Body:   []rules.Point{},
				Health: 100,
			},
		},
		Food: []rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}},
	},
	[]SnakeMove{
		{ID: "one", Move: rules.MoveUp},
		{ID: "two", Move: rules.MoveDown},
	},
	rules.ErrorZeroLengthSnake,
	nil,
}

// Checks a basic state where a snake moves, eats and grows
var standardCaseMoveEatAndGrow = gameTestCase{
	"Standard Case Move Eat and Grow",
	&rules.BoardState{
		Width:  10,
		Height: 10,
		Snakes: []rules.Snake{
			{
				ID:     "one",
				Body:   []rules.Point{{X: 1, Y: 1}, {X: 1, Y: 2}},
				Health: 100,
			},
			{
				ID:     "two",
				Body:   []rules.Point{{X: 3, Y: 4}, {X: 3, Y: 3}},
				Health: 100,
			},
			{
				ID:              "three",
				Body:            []rules.Point{},
				Health:          100,
				EliminatedCause: rules.EliminatedByOutOfBounds,
			},
		},
		Food: []rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}},
	},
	[]SnakeMove{
		{ID: "one", Move: rules.MoveDown},
		{ID: "two", Move: rules.MoveUp},
		{ID: "three", Move: rules.MoveLeft}, // Should be ignored
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
			{
				ID:     "two",
				Body:   []rules.Point{{X: 3, Y: 5}, {X: 3, Y: 4}},
				Health: 99,
			},
			{
				ID:              "three",
				Body:            []rules.Point{},
				Health:          100,
				EliminatedCause: rules.EliminatedByOutOfBounds,
			},
		},
		Food: []rules.Point{{X: 0, Y: 0}},
	},
}

// Checks a basic state where two snakes of equal sizes collide, and both should
// be eliminated as a result.
var standardMoveAndCollideMAD = gameTestCase{
	"Standard Case Move and Collide",
	&rules.BoardState{
		Turn:   0,
		Width:  10,
		Height: 10,
		Snakes: []rules.Snake{
			{
				ID:     "one",
				Body:   []rules.Point{{X: 1, Y: 1}, {X: 2, Y: 1}},
				Health: 99,
			},
			{
				ID:     "two",
				Body:   []rules.Point{{X: 1, Y: 2}, {X: 2, Y: 2}},
				Health: 99,
			},
		},
		Food: []rules.Point{},
	},
	[]SnakeMove{
		{ID: "one", Move: rules.MoveUp},
		{ID: "two", Move: rules.MoveDown},
	},
	nil,
	&rules.BoardState{
		Width:  10,
		Height: 10,
		Snakes: []rules.Snake{
			{
				ID:               "one",
				Body:             []rules.Point{{X: 1, Y: 2}, {X: 1, Y: 1}},
				Health:           98,
				EliminatedCause:  rules.EliminatedByCollision,
				EliminatedBy:     "two",
				EliminatedOnTurn: 1,
			},
			{
				ID:               "two",
				Body:             []rules.Point{{X: 1, Y: 1}, {X: 1, Y: 2}},
				Health:           98,
				EliminatedCause:  rules.EliminatedByCollision,
				EliminatedBy:     "one",
				EliminatedOnTurn: 1,
			},
		},
		Food: []rules.Point{},
	},
}

func TestStandardCreateNextBoardState(t *testing.T) {
	cases := []gameTestCase{
		standardCaseErrNoMoveFound,
		standardCaseErrZeroLengthSnake,
		standardCaseMoveEatAndGrow,
		standardMoveAndCollideMAD,
	}
	r := getStandardRuleset(settings.Settings{})
	for _, gc := range cases {
		// test a RulesBuilder constructed instance
		gc.requireValidNextState(t, r)
		// also test a pipeline with the same settings
		gc.requireValidNextState(t, NewRulesetBuilder().PipelineRuleset(rules.GameTypeStandard, NewPipeline(standardRulesetStages...)))
	}
}

func TestEatingOnLastMove(t *testing.T) {
	// We want to specifically ensure that snakes eating food on their last turn
	// survive. It used to be that this wasn't the case, and snakes were eliminated
	// if they moved onto food with their final move. This behaviour wasn't "wrong" or incorrect,
	// it just was less fun to watch. So let's ensure we're always giving snakes every possible
	// changes to reach food before eliminating them.
	tests := []struct {
		prevState     *rules.BoardState
		moves         []SnakeMove
		expectedError error
		expectedState *rules.BoardState
	}{
		{
			&rules.BoardState{
				Width:  10,
				Height: 10,
				Snakes: []rules.Snake{
					{
						ID:     "one",
						Body:   []rules.Point{{X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 0}},
						Health: 1,
					},
					{
						ID:     "two",
						Body:   []rules.Point{{X: 3, Y: 2}, {X: 3, Y: 3}, {X: 3, Y: 4}},
						Health: 1,
					},
				},
				Food: []rules.Point{{X: 0, Y: 3}, {X: 9, Y: 9}},
			},
			[]SnakeMove{
				{ID: "one", Move: rules.MoveUp},
				{ID: "two", Move: rules.MoveDown},
			},
			nil,
			&rules.BoardState{
				Width:  10,
				Height: 10,
				Snakes: []rules.Snake{
					{
						ID:     "one",
						Body:   []rules.Point{{X: 0, Y: 3}, {X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 1}},
						Health: 100,
					},
					{
						ID:               "two",
						Body:             []rules.Point{{X: 3, Y: 1}, {X: 3, Y: 2}, {X: 3, Y: 3}},
						Health:           0,
						EliminatedCause:  rules.EliminatedByOutOfHealth,
						EliminatedOnTurn: 1,
					},
				},
				Food: []rules.Point{{X: 9, Y: 9}},
			},
		},
	}

	rand.Seed(0) // Seed with a value that will reliably not spawn food
	r := getStandardRuleset(settings.Settings{})
	for _, test := range tests {
		_, nextState, err := r.Execute(test.prevState, test.moves)
		require.Equal(t, err, test.expectedError)
		if test.expectedState != nil {
			require.Equal(t, test.expectedState.Width, nextState.Width)
			require.Equal(t, test.expectedState.Height, nextState.Height)
			require.Equal(t, test.expectedState.Food, nextState.Food)
			require.Equal(t, test.expectedState.Snakes, nextState.Snakes)
		}
	}
}

func TestHeadToHeadOnFood(t *testing.T) {
	// We want to specifically ensure that snakes that collide head-to-head
	// on top of food successfully remove the food - that's the core behaviour this test
	// is enforcing. There's a known side effect of this though, in that both snakes will
	// have eaten prior to being evaluated on the head-to-head (+1 length, full health).
	// We're okay with that since it does not impact the result of the head-to-head,
	// however that behaviour could change in the future and this test could be updated.
	tests := []struct {
		prevState     *rules.BoardState
		moves         []SnakeMove
		expectedError error
		expectedState *rules.BoardState
	}{
		{
			&rules.BoardState{
				Turn:   41,
				Width:  10,
				Height: 10,
				Snakes: []rules.Snake{
					{
						ID:     "one",
						Body:   []rules.Point{{X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 0}},
						Health: 10,
					},
					{
						ID:     "two",
						Body:   []rules.Point{{X: 0, Y: 4}, {X: 0, Y: 5}, {X: 0, Y: 6}},
						Health: 10,
					},
				},
				Food: []rules.Point{{X: 0, Y: 3}, {X: 9, Y: 9}},
			},
			[]SnakeMove{
				{ID: "one", Move: rules.MoveUp},
				{ID: "two", Move: rules.MoveDown},
			},
			nil,
			&rules.BoardState{
				Width:  10,
				Height: 10,
				Snakes: []rules.Snake{
					{
						ID:               "one",
						Body:             []rules.Point{{X: 0, Y: 3}, {X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 1}},
						Health:           100,
						EliminatedCause:  rules.EliminatedByHeadToHeadCollision,
						EliminatedBy:     "two",
						EliminatedOnTurn: 42,
					},
					{
						ID:               "two",
						Body:             []rules.Point{{X: 0, Y: 3}, {X: 0, Y: 4}, {X: 0, Y: 5}, {X: 0, Y: 5}},
						Health:           100,
						EliminatedCause:  rules.EliminatedByHeadToHeadCollision,
						EliminatedBy:     "one",
						EliminatedOnTurn: 42,
					},
				},
				Food: []rules.Point{{X: 9, Y: 9}},
			},
		},
		{
			&rules.BoardState{
				Turn:   41,
				Width:  10,
				Height: 10,
				Snakes: []rules.Snake{
					{
						ID:     "one",
						Body:   []rules.Point{{X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 0}},
						Health: 10,
					},
					{
						ID:     "two",
						Body:   []rules.Point{{X: 0, Y: 4}, {X: 0, Y: 5}, {X: 0, Y: 6}, {X: 0, Y: 7}},
						Health: 10,
					},
				},
				Food: []rules.Point{{X: 0, Y: 3}, {X: 9, Y: 9}},
			},
			[]SnakeMove{
				{ID: "one", Move: rules.MoveUp},
				{ID: "two", Move: rules.MoveDown},
			},
			nil,
			&rules.BoardState{
				Width:  10,
				Height: 10,
				Snakes: []rules.Snake{
					{
						ID:               "one",
						Body:             []rules.Point{{X: 0, Y: 3}, {X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 1}},
						Health:           100,
						EliminatedCause:  rules.EliminatedByHeadToHeadCollision,
						EliminatedBy:     "two",
						EliminatedOnTurn: 42,
					},
					{
						ID:     "two",
						Body:   []rules.Point{{X: 0, Y: 3}, {X: 0, Y: 4}, {X: 0, Y: 5}, {X: 0, Y: 6}, {X: 0, Y: 6}},
						Health: 100,
					},
				},
				Food: []rules.Point{{X: 9, Y: 9}},
			},
		},
	}

	rand.Seed(0) // Seed with a value that will reliably not spawn food
	r := getStandardRuleset(settings.Settings{})
	for _, test := range tests {
		_, nextState, err := r.Execute(test.prevState, test.moves)
		require.Equal(t, test.expectedError, err)
		if test.expectedState != nil {
			require.Equal(t, test.expectedState.Width, nextState.Width)
			require.Equal(t, test.expectedState.Height, nextState.Height)
			require.Equal(t, test.expectedState.Food, nextState.Food)
			require.Equal(t, test.expectedState.Snakes, nextState.Snakes)
		}
	}
}

func TestRegressionIssue19(t *testing.T) {
	// Eliminated snakes passed to CreateNextBoardState should not impact next game state
	tests := []struct {
		prevState     *rules.BoardState
		moves         []SnakeMove
		expectedError error
		expectedState *rules.BoardState
	}{
		{
			&rules.BoardState{
				Width:  10,
				Height: 10,
				Snakes: []rules.Snake{
					{
						ID:     "one",
						Body:   []rules.Point{{X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 0}},
						Health: 100,
					},
					{
						ID:     "two",
						Body:   []rules.Point{{X: 0, Y: 5}, {X: 0, Y: 6}, {X: 0, Y: 7}},
						Health: 100,
					},
					{
						ID:              "eliminated",
						Body:            []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}, {X: 0, Y: 4}, {X: 0, Y: 5}, {X: 0, Y: 6}},
						Health:          0,
						EliminatedCause: rules.EliminatedByOutOfHealth,
					},
				},
				Food: []rules.Point{{X: 9, Y: 9}},
			},
			[]SnakeMove{
				{ID: "one", Move: rules.MoveUp},
				{ID: "two", Move: rules.MoveDown},
			},
			nil,
			&rules.BoardState{
				Width:  10,
				Height: 10,
				Snakes: []rules.Snake{
					{
						ID:     "one",
						Body:   []rules.Point{{X: 0, Y: 3}, {X: 0, Y: 2}, {X: 0, Y: 1}},
						Health: 99,
					},
					{
						ID:     "two",
						Body:   []rules.Point{{X: 0, Y: 4}, {X: 0, Y: 5}, {X: 0, Y: 6}},
						Health: 99,
					},
					{
						ID:              "eliminated",
						Body:            []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}, {X: 0, Y: 4}, {X: 0, Y: 5}, {X: 0, Y: 6}},
						Health:          0,
						EliminatedCause: rules.EliminatedByOutOfHealth,
					},
				},
				Food: []rules.Point{{X: 9, Y: 9}},
			},
		},
	}

	rand.Seed(0) // Seed with a value that will reliably not spawn food
	r := getStandardRuleset(settings.Settings{})
	for _, test := range tests {
		_, nextState, err := r.Execute(test.prevState, test.moves)
		require.Equal(t, err, test.expectedError)
		if test.expectedState != nil {
			require.Equal(t, test.expectedState.Width, nextState.Width)
			require.Equal(t, test.expectedState.Height, nextState.Height)
			require.Equal(t, test.expectedState.Food, nextState.Food)
			require.Equal(t, test.expectedState.Snakes, nextState.Snakes)
		}
	}

}

func TestMoveSnakes(t *testing.T) {
	b := &rules.BoardState{
		Snakes: []rules.Snake{
			{
				ID:     "one",
				Body:   []rules.Point{{X: 10, Y: 110}, {X: 11, Y: 110}},
				Health: 111111,
			},
			{
				ID:     "two",
				Body:   []rules.Point{{X: 23, Y: 220}, {X: 22, Y: 220}, {X: 21, Y: 220}, {X: 20, Y: 220}},
				Health: 222222,
			},
			{
				ID:              "three",
				Body:            []rules.Point{{X: 0, Y: 0}},
				Health:          1,
				EliminatedCause: rules.EliminatedByOutOfBounds,
			},
		},
	}

	tests := []struct {
		MoveOne       string
		ExpectedOne   []rules.Point
		MoveTwo       string
		ExpectedTwo   []rules.Point
		MoveThree     string
		ExpectedThree []rules.Point
	}{
		{
			rules.MoveDown, []rules.Point{{X: 10, Y: 109}, {X: 10, Y: 110}},
			rules.MoveUp, []rules.Point{{X: 23, Y: 221}, {X: 23, Y: 220}, {X: 22, Y: 220}, {X: 21, Y: 220}},
			rules.MoveDown, []rules.Point{{X: 0, Y: 0}},
		},
		{
			rules.MoveRight, []rules.Point{{X: 11, Y: 109}, {X: 10, Y: 109}},
			rules.MoveLeft, []rules.Point{{X: 22, Y: 221}, {X: 23, Y: 221}, {X: 23, Y: 220}, {X: 22, Y: 220}},
			rules.MoveDown, []rules.Point{{X: 0, Y: 0}},
		},
		{
			rules.MoveRight, []rules.Point{{X: 12, Y: 109}, {X: 11, Y: 109}},
			rules.MoveLeft, []rules.Point{{X: 21, Y: 221}, {X: 22, Y: 221}, {X: 23, Y: 221}, {X: 23, Y: 220}},
			rules.MoveDown, []rules.Point{{X: 0, Y: 0}},
		},
		{
			rules.MoveRight, []rules.Point{{X: 13, Y: 109}, {X: 12, Y: 109}},
			rules.MoveLeft, []rules.Point{{X: 20, Y: 221}, {X: 21, Y: 221}, {X: 22, Y: 221}, {X: 23, Y: 221}},
			rules.MoveDown, []rules.Point{{X: 0, Y: 0}},
		},
		{
			rules.MoveDown, []rules.Point{{X: 13, Y: 108}, {X: 13, Y: 109}},
			rules.MoveUp, []rules.Point{{X: 20, Y: 222}, {X: 20, Y: 221}, {X: 21, Y: 221}, {X: 22, Y: 221}},
			rules.MoveDown, []rules.Point{{X: 0, Y: 0}},
		},
	}

	r := getStandardRuleset(settings.Settings{})
	for _, test := range tests {
		moves := []SnakeMove{
			{ID: "one", Move: test.MoveOne},
			{ID: "two", Move: test.MoveTwo},
			{ID: "three", Move: test.MoveThree},
		}
		_, err := MoveSnakesStandard(b, r.Settings(), moves)

		require.NoError(t, err)
		require.Len(t, b.Snakes, 3)

		require.Equal(t, 111111, b.Snakes[0].Health)
		require.Equal(t, 222222, b.Snakes[1].Health)
		require.Equal(t, 1, b.Snakes[2].Health)

		require.Len(t, b.Snakes[0].Body, 2)
		require.Len(t, b.Snakes[1].Body, 4)
		require.Len(t, b.Snakes[2].Body, 1)

		require.Equal(t, len(b.Snakes[0].Body), len(test.ExpectedOne))
		for i, e := range test.ExpectedOne {
			require.Equal(t, e, b.Snakes[0].Body[i])
		}
		require.Equal(t, len(b.Snakes[1].Body), len(test.ExpectedTwo))
		for i, e := range test.ExpectedTwo {
			require.Equal(t, e, b.Snakes[1].Body[i])
		}
		require.Equal(t, len(b.Snakes[2].Body), len(test.ExpectedThree))
		for i, e := range test.ExpectedThree {
			require.Equal(t, e, b.Snakes[2].Body[i])
		}
	}
}

func TestMoveSnakesWrongID(t *testing.T) {
	b := &rules.BoardState{
		Snakes: []rules.Snake{
			{
				ID:   "one",
				Body: []rules.Point{{X: 1, Y: 1}},
			},
		},
	}
	moves := []SnakeMove{
		{
			ID:   "not found",
			Move: rules.MoveUp,
		},
	}

	r := getStandardRuleset(settings.Settings{})
	_, err := MoveSnakesStandard(b, r.Settings(), moves)
	require.Equal(t, rules.ErrorNoMoveFound, err)
}

func TestMoveSnakesNotEnoughMoves(t *testing.T) {
	b := &rules.BoardState{
		Snakes: []rules.Snake{
			{
				ID:   "one",
				Body: []rules.Point{{X: 1, Y: 1}},
			},
			{
				ID:   "two",
				Body: []rules.Point{{X: 2, Y: 2}},
			},
		},
	}
	moves := []SnakeMove{
		{
			ID:   "two",
			Move: rules.MoveUp,
		},
	}

	r := getStandardRuleset(settings.Settings{})
	_, err := MoveSnakesStandard(b, r.Settings(), moves)
	require.Equal(t, rules.ErrorNoMoveFound, err)
}

func TestMoveSnakesExtraMovesIgnored(t *testing.T) {
	b := &rules.BoardState{
		Snakes: []rules.Snake{
			{
				ID:   "one",
				Body: []rules.Point{{X: 1, Y: 1}},
			},
		},
	}
	moves := []SnakeMove{
		{
			ID:   "one",
			Move: rules.MoveDown,
		},
		{
			ID:   "two",
			Move: rules.MoveLeft,
		},
	}

	r := getStandardRuleset(settings.Settings{})
	_, err := MoveSnakesStandard(b, r.Settings(), moves)
	require.NoError(t, err)
	require.Equal(t, []rules.Point{{X: 1, Y: 0}}, b.Snakes[0].Body)
}

func TestMoveSnakesDefault(t *testing.T) {
	tests := []struct {
		Body     []rules.Point
		Move     string
		Expected []rules.Point
	}{
		{
			Body:     []rules.Point{{X: 0, Y: 0}},
			Move:     "invalid",
			Expected: []rules.Point{{X: 0, Y: 1}},
		},
		{
			Body:     []rules.Point{{X: 5, Y: 5}, {X: 5, Y: 5}},
			Move:     "",
			Expected: []rules.Point{{X: 5, Y: 6}, {X: 5, Y: 5}},
		},
		{
			Body:     []rules.Point{{X: 5, Y: 5}, {X: 5, Y: 4}},
			Expected: []rules.Point{{X: 5, Y: 6}, {X: 5, Y: 5}},
		},
		{
			Body:     []rules.Point{{X: 5, Y: 4}, {X: 5, Y: 5}},
			Expected: []rules.Point{{X: 5, Y: 3}, {X: 5, Y: 4}},
		},
		{
			Body:     []rules.Point{{X: 5, Y: 4}, {X: 5, Y: 5}},
			Expected: []rules.Point{{X: 5, Y: 3}, {X: 5, Y: 4}},
		},
		{
			Body:     []rules.Point{{X: 4, Y: 5}, {X: 5, Y: 5}},
			Expected: []rules.Point{{X: 3, Y: 5}, {X: 4, Y: 5}},
		},
		{
			Body:     []rules.Point{{X: 5, Y: 5}, {X: 4, Y: 5}},
			Expected: []rules.Point{{X: 6, Y: 5}, {X: 5, Y: 5}},
		},
	}

	r := getStandardRuleset(settings.Settings{})
	for _, test := range tests {
		b := &rules.BoardState{
			Snakes: []rules.Snake{
				{ID: "one", Body: test.Body},
			},
		}
		moves := []SnakeMove{{ID: "one", Move: test.Move}}

		_, err := MoveSnakesStandard(b, r.Settings(), moves)
		require.NoError(t, err)
		require.Len(t, b.Snakes, 1)
		require.Equal(t, len(test.Body), len(b.Snakes[0].Body))
		require.Equal(t, len(test.Expected), len(b.Snakes[0].Body))
		for i, e := range test.Expected {
			require.Equal(t, e, b.Snakes[0].Body[i])
		}
	}
}

func TestGetDefaultMove(t *testing.T) {
	tests := []struct {
		SnakeBody    []rules.Point
		ExpectedMove string
	}{
		// Default is always up
		{
			SnakeBody:    []rules.Point{},
			ExpectedMove: rules.MoveUp,
		},
		{
			SnakeBody:    []rules.Point{{X: 0, Y: 0}},
			ExpectedMove: rules.MoveUp,
		},
		{
			SnakeBody:    []rules.Point{{X: -1, Y: -1}},
			ExpectedMove: rules.MoveUp,
		},
		// Stacked (fallback to default)
		{
			SnakeBody:    []rules.Point{{X: 2, Y: 2}, {X: 2, Y: 2}},
			ExpectedMove: rules.MoveUp,
		},
		// Neck next to head
		{
			SnakeBody:    []rules.Point{{X: 2, Y: 2}, {X: 2, Y: 1}},
			ExpectedMove: rules.MoveUp,
		},
		{
			SnakeBody:    []rules.Point{{X: 2, Y: 2}, {X: 2, Y: 3}},
			ExpectedMove: rules.MoveDown,
		},
		{
			SnakeBody:    []rules.Point{{X: 2, Y: 2}, {X: 1, Y: 2}},
			ExpectedMove: rules.MoveRight,
		},
		{
			SnakeBody:    []rules.Point{{X: 2, Y: 2}, {X: 3, Y: 2}},
			ExpectedMove: rules.MoveLeft,
		},
		// Board wrap cases
		{
			SnakeBody:    []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 2}},
			ExpectedMove: rules.MoveUp,
		},
		{
			SnakeBody:    []rules.Point{{X: 0, Y: 0}, {X: 2, Y: 0}},
			ExpectedMove: rules.MoveRight,
		},
		{
			SnakeBody:    []rules.Point{{X: 0, Y: 2}, {X: 0, Y: 0}},
			ExpectedMove: rules.MoveDown,
		},
		{
			SnakeBody:    []rules.Point{{X: 2, Y: 0}, {X: 0, Y: 0}},
			ExpectedMove: rules.MoveLeft,
		},
	}

	for _, test := range tests {
		actualMove := getDefaultMove(test.SnakeBody)
		require.Equal(t, test.ExpectedMove, actualMove)
	}
}

func TestReduceSnakeHealth(t *testing.T) {
	b := &rules.BoardState{
		Snakes: []rules.Snake{
			{
				Body:   []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}},
				Health: 99,
			},
			{
				Body:   []rules.Point{{X: 5, Y: 8}, {X: 6, Y: 8}, {X: 7, Y: 8}},
				Health: 2,
			},
			{
				Body:            []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}},
				Health:          50,
				EliminatedCause: rules.EliminatedByCollision,
			},
		},
	}

	r := getStandardRuleset(settings.Settings{})
	_, err := ReduceSnakeHealthStandard(b, r.Settings(), mockSnakeMoves())
	require.NoError(t, err)
	require.Equal(t, b.Snakes[0].Health, 98)
	require.Equal(t, b.Snakes[1].Health, 1)
	require.Equal(t, b.Snakes[2].Health, 50)

	_, err = ReduceSnakeHealthStandard(b, r.Settings(), mockSnakeMoves())
	require.NoError(t, err)
	require.Equal(t, b.Snakes[0].Health, 97)
	require.Equal(t, b.Snakes[1].Health, 0)
	require.Equal(t, b.Snakes[2].Health, 50)

	_, err = ReduceSnakeHealthStandard(b, r.Settings(), mockSnakeMoves())
	require.NoError(t, err)
	require.Equal(t, b.Snakes[0].Health, 96)
	require.Equal(t, b.Snakes[1].Health, -1)
	require.Equal(t, b.Snakes[2].Health, 50)

	_, err = ReduceSnakeHealthStandard(b, r.Settings(), mockSnakeMoves())
	require.NoError(t, err)
	require.Equal(t, b.Snakes[0].Health, 95)
	require.Equal(t, b.Snakes[1].Health, -2)
	require.Equal(t, b.Snakes[2].Health, 50)
}

func TestSnakeIsOutOfHealth(t *testing.T) {
	tests := []struct {
		Health   int
		Expected bool
	}{
		{Health: math.MinInt, Expected: true},
		{Health: -10, Expected: true},
		{Health: -2, Expected: true},
		{Health: -1, Expected: true},
		{Health: 0, Expected: true},
		{Health: 1, Expected: false},
		{Health: 2, Expected: false},
		{Health: 10, Expected: false},
		{Health: math.MaxInt, Expected: false},
	}

	for _, test := range tests {
		s := &rules.Snake{Health: test.Health}
		require.Equal(t, test.Expected, snakeIsOutOfHealth(s), "Health: %+v", test.Health)
	}
}

func TestSnakeIsOutOfBounds(t *testing.T) {
	boardWidth := 10
	boardHeight := 100

	tests := []struct {
		Point    rules.Point
		Expected bool
	}{
		{rules.Point{X: math.MinInt, Y: math.MinInt}, true},
		{rules.Point{X: math.MinInt, Y: 0}, true},
		{rules.Point{X: 0, Y: math.MinInt}, true},
		{rules.Point{X: -1, Y: -1}, true},
		{rules.Point{X: -1, Y: 0}, true},
		{rules.Point{X: 0, Y: -1}, true},
		{rules.Point{X: 0, Y: 0}, false},
		{rules.Point{X: 1, Y: 0}, false},
		{rules.Point{X: 0, Y: 1}, false},
		{rules.Point{X: 1, Y: 1}, false},
		{rules.Point{X: 9, Y: 9}, false},
		{rules.Point{X: 9, Y: 10}, false},
		{rules.Point{X: 9, Y: 11}, false},
		{rules.Point{X: 10, Y: 9}, true},
		{rules.Point{X: 10, Y: 10}, true},
		{rules.Point{X: 10, Y: 11}, true},
		{rules.Point{X: 11, Y: 9}, true},
		{rules.Point{X: 11, Y: 10}, true},
		{rules.Point{X: 11, Y: 11}, true},
		{rules.Point{X: math.MaxInt, Y: 11}, true},
		{rules.Point{X: 9, Y: 99}, false},
		{rules.Point{X: 9, Y: 100}, true},
		{rules.Point{X: 9, Y: 101}, true},
		{rules.Point{X: 9, Y: math.MaxInt}, true},
		{rules.Point{X: math.MaxInt, Y: math.MaxInt}, true},
	}

	for _, test := range tests {
		// Test with point as head
		s := rules.Snake{Body: []rules.Point{test.Point}}
		require.Equal(t, test.Expected, snakeIsOutOfBounds(&s, boardWidth, boardHeight), "Head%+v", test.Point)
		// Test with point as body
		s = rules.Snake{Body: []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 0}, test.Point}}
		require.Equal(t, test.Expected, snakeIsOutOfBounds(&s, boardWidth, boardHeight), "Body%+v", test.Point)
	}
}

func TestSnakeHasBodyCollidedSelf(t *testing.T) {
	tests := []struct {
		Body     []rules.Point
		Expected bool
	}{
		{[]rules.Point{{X: 1, Y: 1}}, false},
		// Self stacks should self collide
		// (we rely on snakes moving before we check self-collision on turn one)
		{[]rules.Point{{X: 2, Y: 2}, {X: 2, Y: 2}}, true},
		{[]rules.Point{{X: 3, Y: 3}, {X: 3, Y: 3}, {X: 3, Y: 3}}, true},
		{[]rules.Point{{X: 5, Y: 5}, {X: 5, Y: 5}, {X: 5, Y: 5}, {X: 5, Y: 5}, {X: 5, Y: 5}}, true},
		// Non-collision cases
		{[]rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 0}}, false},
		{[]rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}}, false},
		{[]rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}}, false},
		{[]rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 4, Y: 0}}, false},
		{[]rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}, false},
		{[]rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 2}}, false},
		{[]rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}, {X: 0, Y: 4}}, false},
		// Collision cases
		{[]rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}}, true},
		{[]rules.Point{{X: 0, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}}, true},
		{[]rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1}, {X: 0, Y: 0}}, true},
		{[]rules.Point{{X: 4, Y: 4}, {X: 3, Y: 4}, {X: 3, Y: 3}, {X: 4, Y: 4}, {X: 4, Y: 4}}, true},
		{[]rules.Point{{X: 3, Y: 3}, {X: 3, Y: 4}, {X: 3, Y: 3}, {X: 4, Y: 4}, {X: 4, Y: 5}}, true},
	}

	for _, test := range tests {
		s := rules.Snake{Body: test.Body}
		require.Equal(t, test.Expected, snakeHasBodyCollided(&s, &s), "Body%q", s.Body)
	}
}

func TestSnakeHasBodyCollidedOther(t *testing.T) {
	tests := []struct {
		SnakeBody []rules.Point
		OtherBody []rules.Point
		Expected  bool
	}{
		{
			// Just heads
			[]rules.Point{{X: 0, Y: 0}},
			[]rules.Point{{X: 1, Y: 1}},
			false,
		},
		{
			// Head-to-heads are not considered in body collisions
			[]rules.Point{{X: 0, Y: 0}},
			[]rules.Point{{X: 0, Y: 0}},
			false,
		},
		{
			// Stacked bodies
			[]rules.Point{{X: 0, Y: 0}},
			[]rules.Point{{X: 0, Y: 0}, {X: 0, Y: 0}},
			true,
		},
		{
			// Separate stacked bodies
			[]rules.Point{{X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}},
			[]rules.Point{{X: 1, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 1}},
			false,
		},
		{
			// Stacked bodies, separated heads
			[]rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 0}},
			[]rules.Point{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 0}},
			false,
		},
		{
			// Mid-snake collision
			[]rules.Point{{X: 1, Y: 1}},
			[]rules.Point{{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}},
			true,
		},
	}

	for _, test := range tests {
		s := &rules.Snake{Body: test.SnakeBody}
		o := &rules.Snake{Body: test.OtherBody}
		require.Equal(t, test.Expected, snakeHasBodyCollided(s, o), "Snake%q Other%q", s.Body, o.Body)
	}
}

func TestSnakeHasLostHeadToHead(t *testing.T) {
	tests := []struct {
		SnakeBody        []rules.Point
		OtherBody        []rules.Point
		Expected         bool
		ExpectedOpposite bool
	}{
		{
			// Just heads
			[]rules.Point{{X: 0, Y: 0}},
			[]rules.Point{{X: 1, Y: 1}},
			false, false,
		},
		{
			// Just heads colliding
			[]rules.Point{{X: 0, Y: 0}},
			[]rules.Point{{X: 0, Y: 0}},
			true, true,
		},
		{
			// One snake larger
			[]rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
			[]rules.Point{{X: 0, Y: 0}},
			false, true,
		},
		{
			// Other snake equal
			[]rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
			[]rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}},
			true, true,
		},
		{
			// Other snake longer
			[]rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
			[]rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}},
			true, false,
		},
		{
			// Body collision
			[]rules.Point{{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}},
			[]rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}},
			false, false,
		},
		{
			// Separate stacked bodies, head collision
			[]rules.Point{{X: 3, Y: 10}, {X: 2, Y: 10}, {X: 2, Y: 10}},
			[]rules.Point{{X: 3, Y: 10}, {X: 4, Y: 10}, {X: 4, Y: 10}},
			true, true,
		},
		{
			// Separate stacked bodies, head collision
			[]rules.Point{{X: 10, Y: 3}, {X: 10, Y: 2}, {X: 10, Y: 1}, {X: 10, Y: 0}},
			[]rules.Point{{X: 10, Y: 3}, {X: 10, Y: 4}, {X: 10, Y: 5}},
			false, true,
		},
	}

	for _, test := range tests {
		s := rules.Snake{Body: test.SnakeBody}
		o := rules.Snake{Body: test.OtherBody}
		require.Equal(t, test.Expected, snakeHasLostHeadToHead(&s, &o), "Snake%q Other%q", s.Body, o.Body)
		require.Equal(t, test.ExpectedOpposite, snakeHasLostHeadToHead(&o, &s), "Snake%q Other%q", s.Body, o.Body)
	}

}

func TestMaybeEliminateSnakes(t *testing.T) {
	tests := []struct {
		Name                     string
		Snakes                   []rules.Snake
		ExpectedEliminatedCauses []string
		ExpectedEliminatedBy     []string
		Err                      error
	}{
		{
			"Empty",
			[]rules.Snake{},
			[]string{},
			[]string{},
			nil,
		},
		{
			"Zero Snake",
			[]rules.Snake{
				{},
			},
			[]string{rules.NotEliminated},
			[]string{""},
			rules.ErrorZeroLengthSnake,
		},
		{
			"Single Starvation",
			[]rules.Snake{
				{ID: "1", Body: []rules.Point{{X: 1, Y: 1}}},
			},
			[]string{rules.EliminatedByOutOfHealth},
			[]string{""},
			nil,
		},
		{
			"Not Eliminated",
			[]rules.Snake{
				{ID: "1", Health: 1, Body: []rules.Point{{X: 1, Y: 1}}},
			},
			[]string{rules.NotEliminated},
			[]string{""},
			nil,
		},
		{
			"Out of Bounds",
			[]rules.Snake{
				{ID: "1", Health: 1, Body: []rules.Point{{X: -1, Y: 1}}},
			},
			[]string{rules.EliminatedByOutOfBounds},
			[]string{""},
			nil,
		},
		{
			"Self Collision",
			[]rules.Snake{
				{ID: "1", Health: 1, Body: []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 0}}},
			},
			[]string{rules.EliminatedBySelfCollision},
			[]string{"1"},
			nil,
		},
		{
			"Multiple Separate Deaths",
			[]rules.Snake{
				{ID: "1", Health: 1, Body: []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 0}}},
				{ID: "2", Health: 1, Body: []rules.Point{{X: -1, Y: 1}}},
			},
			[]string{
				rules.EliminatedBySelfCollision,
				rules.EliminatedByOutOfBounds},
			[]string{"1", ""},
			nil,
		},
		{
			"Other Collision",
			[]rules.Snake{
				{ID: "1", Health: 1, Body: []rules.Point{{X: 0, Y: 2}, {X: 0, Y: 3}, {X: 0, Y: 4}}},
				{ID: "2", Health: 1, Body: []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}},
			},
			[]string{
				rules.EliminatedByCollision,
				rules.NotEliminated},
			[]string{"2", ""},
			nil,
		},
		{
			"All Eliminated Head 2 Head",
			[]rules.Snake{
				{ID: "1", Health: 1, Body: []rules.Point{{X: 1, Y: 1}}},
				{ID: "2", Health: 1, Body: []rules.Point{{X: 1, Y: 1}}},
				{ID: "3", Health: 1, Body: []rules.Point{{X: 1, Y: 1}}},
			},
			[]string{
				rules.EliminatedByHeadToHeadCollision,
				rules.EliminatedByHeadToHeadCollision,
				rules.EliminatedByHeadToHeadCollision,
			},
			[]string{"2", "1", "1"},
			nil,
		},
		{
			"One Snake wins Head 2 Head",
			[]rules.Snake{
				{ID: "1", Health: 1, Body: []rules.Point{{X: 1, Y: 1}, {X: 0, Y: 1}}},
				{ID: "2", Health: 1, Body: []rules.Point{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 1, Y: 3}}},
				{ID: "3", Health: 1, Body: []rules.Point{{X: 1, Y: 1}}},
			},
			[]string{
				rules.EliminatedByHeadToHeadCollision,
				rules.NotEliminated,
				rules.EliminatedByHeadToHeadCollision,
			},
			[]string{"2", "", "2"},
			nil,
		},
		{
			"All Snakes Body Eliminated",
			[]rules.Snake{
				{ID: "1", Health: 1, Body: []rules.Point{{X: 4, Y: 4}, {X: 3, Y: 3}}},
				{ID: "2", Health: 1, Body: []rules.Point{{X: 3, Y: 3}, {X: 2, Y: 2}}},
				{ID: "3", Health: 1, Body: []rules.Point{{X: 2, Y: 2}, {X: 1, Y: 1}}},
				{ID: "4", Health: 1, Body: []rules.Point{{X: 1, Y: 1}, {X: 4, Y: 4}}},
				{ID: "5", Health: 1, Body: []rules.Point{{X: 4, Y: 4}}}, // Body collision takes priority
			},
			[]string{
				rules.EliminatedByCollision,
				rules.EliminatedByCollision,
				rules.EliminatedByCollision,
				rules.EliminatedByCollision,
				rules.EliminatedByCollision,
			},
			[]string{"4", "1", "2", "3", "4"},
			nil,
		},
		{
			"All Snakes Eliminated Head 2 Head",
			[]rules.Snake{
				{ID: "1", Health: 1, Body: []rules.Point{{X: 4, Y: 4}, {X: 4, Y: 5}}},
				{ID: "2", Health: 1, Body: []rules.Point{{X: 4, Y: 4}, {X: 4, Y: 3}}},
				{ID: "3", Health: 1, Body: []rules.Point{{X: 4, Y: 4}, {X: 5, Y: 4}}},
				{ID: "4", Health: 1, Body: []rules.Point{{X: 4, Y: 4}, {X: 3, Y: 4}}},
			},
			[]string{
				rules.EliminatedByHeadToHeadCollision,
				rules.EliminatedByHeadToHeadCollision,
				rules.EliminatedByHeadToHeadCollision,
				rules.EliminatedByHeadToHeadCollision,
			},
			[]string{"2", "1", "1", "1"},
			nil,
		},
		{
			"4 Snakes Head 2 Head",
			[]rules.Snake{
				{ID: "1", Health: 1, Body: []rules.Point{{X: 4, Y: 4}, {X: 4, Y: 5}}},
				{ID: "2", Health: 1, Body: []rules.Point{{X: 4, Y: 4}, {X: 4, Y: 3}}},
				{ID: "3", Health: 1, Body: []rules.Point{{X: 4, Y: 4}, {X: 5, Y: 4}, {X: 6, Y: 4}}},
				{ID: "4", Health: 1, Body: []rules.Point{{X: 4, Y: 4}, {X: 3, Y: 4}}},
			},
			[]string{
				rules.EliminatedByHeadToHeadCollision,
				rules.EliminatedByHeadToHeadCollision,
				rules.NotEliminated,
				rules.EliminatedByHeadToHeadCollision,
			},
			[]string{"3", "3", "", "3"},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			b := &rules.BoardState{
				Width:  10,
				Height: 10,
				Snakes: test.Snakes,
			}
			_, err := EliminateSnakesStandard(b, settings.Settings{}, mockSnakeMoves())
			require.Equal(t, test.Err, err)
			for i, snake := range b.Snakes {
				require.Equal(t, test.ExpectedEliminatedCauses[i], snake.EliminatedCause)
				require.Equal(t, test.ExpectedEliminatedBy[i], snake.EliminatedBy)
			}
		})
	}
}

func TestMaybeEliminateSnakesPriority(t *testing.T) {
	tests := []struct {
		Snakes                   []rules.Snake
		ExpectedEliminatedCauses []string
		ExpectedEliminatedBy     []string
	}{
		{
			[]rules.Snake{
				{ID: "1", Health: 0, Body: []rules.Point{{X: -1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}}},
				{ID: "2", Health: 1, Body: []rules.Point{{X: -1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}}},
				{ID: "3", Health: 1, Body: []rules.Point{{X: 1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}}},
				{ID: "4", Health: 1, Body: []rules.Point{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}}},
				{ID: "5", Health: 1, Body: []rules.Point{{X: 2, Y: 2}, {X: 2, Y: 1}, {X: 2, Y: 0}}},
				{ID: "6", Health: 1, Body: []rules.Point{{X: 2, Y: 2}, {X: 2, Y: 3}, {X: 2, Y: 4}, {X: 2, Y: 5}}},
			},
			[]string{
				rules.EliminatedByOutOfHealth,
				rules.EliminatedByOutOfBounds,
				rules.EliminatedBySelfCollision,
				rules.EliminatedByCollision,
				rules.EliminatedByHeadToHeadCollision,
				rules.NotEliminated,
			},
			[]string{"", "", "3", "3", "6", ""},
		},
	}

	r := getStandardRuleset(settings.Settings{})
	for _, test := range tests {
		b := &rules.BoardState{Width: 10, Height: 10, Snakes: test.Snakes}
		_, err := EliminateSnakesStandard(b, r.Settings(), mockSnakeMoves())
		require.NoError(t, err)
		for i, snake := range b.Snakes {
			require.Equal(t, test.ExpectedEliminatedCauses[i], snake.EliminatedCause, snake.ID)
			require.Equal(t, test.ExpectedEliminatedBy[i], snake.EliminatedBy, snake.ID)
		}
	}
}

func TestMaybeFeedSnakes(t *testing.T) {
	tests := []struct {
		Name           string
		Snakes         []rules.Snake
		Food           []rules.Point
		ExpectedSnakes []rules.Snake
		ExpectedFood   []rules.Point
	}{
		{
			Name: "snake not on food",
			Snakes: []rules.Snake{
				{Health: 5, Body: []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}},
			},
			Food: []rules.Point{{X: 3, Y: 3}},
			ExpectedSnakes: []rules.Snake{
				{Health: 5, Body: []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}},
			},
			ExpectedFood: []rules.Point{{X: 3, Y: 3}},
		},
		{
			Name: "snake on food",
			Snakes: []rules.Snake{
				{Health: rules.SnakeMaxHealth - 1, Body: []rules.Point{{X: 2, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 2}, {X: 2, Y: 2}}},
			},
			Food: []rules.Point{{X: 2, Y: 1}},
			ExpectedSnakes: []rules.Snake{
				{Health: rules.SnakeMaxHealth, Body: []rules.Point{{X: 2, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 2, Y: 2}}},
			},
			ExpectedFood: []rules.Point{},
		},
		{
			Name: "food under body",
			Snakes: []rules.Snake{
				{Body: []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}},
			},
			Food: []rules.Point{{X: 0, Y: 1}},
			ExpectedSnakes: []rules.Snake{
				{Body: []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}},
			},
			ExpectedFood: []rules.Point{{X: 0, Y: 1}},
		},
		{
			Name: "snake on food but already eliminated",
			Snakes: []rules.Snake{
				{Body: []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}, EliminatedCause: "EliminatedByOutOfBounds"},
			},
			Food: []rules.Point{{X: 0, Y: 0}},
			ExpectedSnakes: []rules.Snake{
				{Body: []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}},
			},
			ExpectedFood: []rules.Point{{X: 0, Y: 0}},
		},
		{
			Name: "multiple snakes on same food",
			Snakes: []rules.Snake{
				{Health: rules.SnakeMaxHealth, Body: []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}},
				{Health: rules.SnakeMaxHealth - 9, Body: []rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}}},
			},
			Food: []rules.Point{{X: 0, Y: 0}, {X: 4, Y: 4}},
			ExpectedSnakes: []rules.Snake{
				{Health: rules.SnakeMaxHealth, Body: []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 2}}},
				{Health: rules.SnakeMaxHealth, Body: []rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}}},
			},
			ExpectedFood: []rules.Point{{X: 4, Y: 4}},
		},
	}

	r := getStandardRuleset(settings.Settings{})
	for _, test := range tests {
		b := &rules.BoardState{
			Snakes: test.Snakes,
			Food:   test.Food,
		}
		_, err := FeedSnakesStandard(b, r.Settings(), nil)
		require.NoError(t, err, test.Name)
		require.Equal(t, len(test.ExpectedSnakes), len(b.Snakes), test.Name)
		for i := 0; i < len(b.Snakes); i++ {
			require.Equal(t, test.ExpectedSnakes[i].Health, b.Snakes[i].Health, test.Name)
			require.Equal(t, test.ExpectedSnakes[i].Body, b.Snakes[i].Body, test.Name)
		}
		require.Equal(t, test.ExpectedFood, b.Food, test.Name)
	}
}

func TestIsGameOver(t *testing.T) {
	tests := []struct {
		Snakes   []rules.Snake
		Expected bool
	}{
		{[]rules.Snake{}, true},
		{[]rules.Snake{{}}, true},
		{[]rules.Snake{{}, {}}, false},
		{[]rules.Snake{{}, {}, {}, {}, {}}, false},
		{
			[]rules.Snake{
				{EliminatedCause: rules.EliminatedByCollision},
				{EliminatedCause: rules.NotEliminated},
			},
			true,
		},
		{
			[]rules.Snake{
				{EliminatedCause: rules.NotEliminated},
				{EliminatedCause: rules.EliminatedByCollision},
				{EliminatedCause: rules.NotEliminated},
				{EliminatedCause: rules.NotEliminated},
			},
			false,
		},
		{
			[]rules.Snake{
				{EliminatedCause: rules.EliminatedByOutOfBounds},
				{EliminatedCause: rules.EliminatedByOutOfBounds},
				{EliminatedCause: rules.EliminatedByOutOfBounds},
				{EliminatedCause: rules.EliminatedByOutOfBounds},
			},
			true,
		},
		{
			[]rules.Snake{
				{EliminatedCause: rules.EliminatedByOutOfBounds},
				{EliminatedCause: rules.EliminatedByOutOfBounds},
				{EliminatedCause: rules.EliminatedByOutOfBounds},
				{EliminatedCause: rules.NotEliminated},
			},
			true,
		},
		{
			[]rules.Snake{
				{EliminatedCause: rules.EliminatedByOutOfBounds},
				{EliminatedCause: rules.EliminatedByOutOfBounds},
				{EliminatedCause: rules.NotEliminated},
				{EliminatedCause: rules.NotEliminated},
			},
			false,
		},
	}

	r := getStandardRuleset(settings.Settings{})
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
