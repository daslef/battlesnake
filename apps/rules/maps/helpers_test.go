package maps_test

import (
	"errors"
	"testing"

	"rules"
	"rules/maps"
	"github.com/stretchr/testify/require"
)

func TestSetupBoard_NotFound(t *testing.T) {
	_, err := maps.SetupBoard("does_not_exist", rules.Settings{}, 10, 10, []string{})

	require.EqualError(t, err, rules.ErrorMapNotFound.Error())
}

func TestSetupBoard_Error(t *testing.T) {
	testMap := maps.StubMap{
		Id:    t.Name(),
		Error: errors.New("bad map update"),
	}
	maps.TestMap(testMap.ID(), testMap, func() {
		_, err := maps.SetupBoard(testMap.ID(), rules.Settings{}, 10, 10, []string{})
		require.EqualError(t, err, "bad map update")
	})
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

	maps.TestMap(testMap.ID(), testMap, func() {
		boardState, err := maps.SetupBoard(testMap.ID(), rules.Settings{}, 10, 10, []string{"1", "2"})

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
	})
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
	maps.TestMap(testMap.ID(), testMap, func() {
		boardState, err := maps.PostUpdateBoard(testMap, previousBoardState, rules.Settings{})

		require.NoError(t, err)

		require.Len(t, boardState.Snakes, 1)

		require.Equal(t, rules.Snake{
			ID:     "1",
			Body:   []rules.Point{{X: 6, Y: 4}, {X: 6, Y: 3}, {X: 6, Y: 2}},
			Health: rules.SnakeMaxHealth,
		}, boardState.Snakes[0])
		require.Equal(t, []rules.Point{{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 5, Y: 3}}, boardState.Food)
	})
}

func TestPlaceFoodFixed(t *testing.T) {
	initialBoardState := rules.NewBoardState(rules.BoardSizeMedium, rules.BoardSizeMedium)
	editor := maps.NewBoardStateEditor(initialBoardState.Clone())

	editor.PlaceSnake("1", []rules.Point{{X: 1, Y: 1}}, 100)
	editor.PlaceSnake("2", []rules.Point{{X: 9, Y: 1}}, 100)
	editor.PlaceSnake("3", []rules.Point{{X: 4, Y: 9}}, 100)
	editor.PlaceSnake("4", []rules.Point{{X: 6, Y: 6}}, 100)

	err := maps.PlaceFoodFixed(rules.MaxRand, initialBoardState, editor)
	require.NoError(t, err)

	food := editor.Food()
	require.Contains(t, food, rules.Point{X: 0, Y: 2})
	require.Contains(t, food, rules.Point{X: 8, Y: 0})
	require.Contains(t, food, rules.Point{X: 3, Y: 10})
	require.Contains(t, food, rules.Point{X: 7, Y: 7})
}
