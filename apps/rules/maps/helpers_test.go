package maps_test

import (
	"errors"
	"testing"

	"rules"
	"rules/maps"
	"rules/settings"

	"github.com/stretchr/testify/require"
)

func TestSetupBoard_Error(t *testing.T) {
	testMap := maps.StubMap{
		Id:    t.Name(),
		Error: errors.New("bad map update"),
	}
	_, err := maps.SetupBoard(testMap, settings.Settings{}, 10, 10, []string{})
	require.EqualError(t, err, "bad map update")
}

func TestSetupBoard(t *testing.T) {
	testMap := maps.StubMap{
		Id: t.Name(),
		SnakePositions: map[string]rules.Point{
			"1": {X: 3, Y: 4},
			"2": {X: 6, Y: 2},
		},
		Food: []rules.Point{
			{X: 1, Y: 1},
			{X: 5, Y: 3},
		},
	}

	boardState, err := maps.SetupBoard(testMap, settings.Settings{}, 10, 10, []string{"1", "2"})

	require.NoError(t, err)

	require.Len(t, boardState.Snakes, 2)

	require.Equal(t, rules.Snake{
		ID:     "1",
		Body:   []rules.Point{{X: 3, Y: 4}, {X: 3, Y: 4}, {X: 3, Y: 4}},
		Health: rules.SnakeMaxHealth,
	}, boardState.Snakes[0])
	require.Equal(t, rules.Snake{
		ID:     "2",
		Body:   []rules.Point{{X: 6, Y: 2}, {X: 6, Y: 2}, {X: 6, Y: 2}},
		Health: rules.SnakeMaxHealth,
	}, boardState.Snakes[1])
	require.Equal(t, []rules.Point{{X: 1, Y: 1}, {X: 5, Y: 3}}, boardState.Food)
}

func TestUpdateBoard(t *testing.T) {
	testMap := maps.StubMap{
		Id: t.Name(),
		SnakePositions: map[string]rules.Point{
			"1": {X: 3, Y: 4},
			"2": {X: 6, Y: 2},
		},
		Food: []rules.Point{
			{X: 1, Y: 1},
			{X: 5, Y: 3},
		},
	}

	previousBoardState := rules.NewBoardState(5, 5).
		WithFood([]rules.Point{{X: 0, Y: 1}}).
		WithSnakes([]rules.Snake{
			{
				ID:     "1",
				Health: 100,
				Body: []rules.Point{
					{X: 6, Y: 4},
					{X: 6, Y: 3},
					{X: 6, Y: 2},
				},
			},
		})

	boardState, err := maps.PostUpdateBoard(testMap, previousBoardState, settings.Settings{})

	require.NoError(t, err)

	require.Len(t, boardState.Snakes, 1)

	require.Equal(t, rules.Snake{
		ID:     "1",
		Body:   []rules.Point{{X: 6, Y: 4}, {X: 6, Y: 3}, {X: 6, Y: 2}},
		Health: rules.SnakeMaxHealth,
	}, boardState.Snakes[0])
	require.Equal(t, []rules.Point{{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 5, Y: 3}}, boardState.Food)

}
