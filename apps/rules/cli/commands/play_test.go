package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"rules"
	"rules/board"
	"rules/client"
	"rules/rulesets"
	"rules/settings"
	"rules/test"

	"github.com/stretchr/testify/require"
)

func buildDefaultGameState() *GameState {
	gameState := &GameState{
		Width:           11,
		Height:          11,
		Names:           nil,
		Timeout:         500,
		GameType:        "standard",
		Seed:            1,
		FoodSpawnChance: 15,
	}

	return gameState
}

func TestGetIndividualBoardStateForSnake(t *testing.T) {
	s1 := rules.Snake{ID: "one", Body: []rules.Point{{X: 3, Y: 3}}}
	s2 := rules.Snake{ID: "two", Body: []rules.Point{{X: 4, Y: 3}}}
	state := rules.NewBoardState(11, 11).
		WithSnakes(
			[]rules.Snake{s1, s2},
		)
	s1State := SnakeState{
		ID:    "one",
		Name:  "ONE",
		URL:   "http://example1.com",
		Head:  "safe",
		Tail:  "curled",
		Color: "#123456",
	}
	s2State := SnakeState{
		ID:    "two",
		Name:  "TWO",
		URL:   "http://example2.com",
		Head:  "silly",
		Tail:  "bolt",
		Color: "#654321",
	}

	gameState := buildDefaultGameState()
	err := gameState.Initialize()
	require.NoError(t, err)
	gameState.snakeStates = map[string]SnakeState{
		s1State.ID: s1State,
		s2State.ID: s2State,
	}

	snakeRequest := gameState.getRequestBodyForSnake(state, s1State)
	requestBody := serialiseSnakeRequest(snakeRequest)

	test.RequireJSONMatchesFixture(t, "testdata/snake_request_body.json", string(requestBody))
}

func TestSettingsRequestSerialization(t *testing.T) {
	s1 := rules.Snake{ID: "one", Body: []rules.Point{{X: 3, Y: 3}}}
	s2 := rules.Snake{ID: "two", Body: []rules.Point{{X: 4, Y: 3}}}
	state := rules.NewBoardState(11, 11).
		WithSnakes([]rules.Snake{s1, s2})
	s1State := SnakeState{
		ID:    "one",
		Name:  "ONE",
		URL:   "http://example1.com",
		Head:  "safe",
		Tail:  "curled",
		Color: "#123456",
	}
	s2State := SnakeState{
		ID:    "two",
		Name:  "TWO",
		URL:   "http://example2.com",
		Head:  "silly",
		Tail:  "bolt",
		Color: "#654321",
	}

	for _, gt := range []string{
		rules.GameTypeStandard, rules.GameTypeSolo,
	} {
		t.Run(gt, func(t *testing.T) {
			gameState := buildDefaultGameState()

			gameState.FoodSpawnChance = 11
			gameState.GameType = gt

			err := gameState.Initialize()
			require.NoError(t, err)
			gameState.snakeStates = map[string]SnakeState{s1State.ID: s1State, s2State.ID: s2State}

			snakeRequest := gameState.getRequestBodyForSnake(state, s1State)
			requestBody := serialiseSnakeRequest(snakeRequest)
			t.Log(string(requestBody))

			test.RequireJSONMatchesFixture(t, fmt.Sprintf("testdata/snake_request_body_%s.json", gt), string(requestBody))
		})
	}
}

func TestConvertRulesSnakes(t *testing.T) {
	tests := []struct {
		name     string
		snakes   []rules.Snake
		state    map[string]SnakeState
		expected []client.Snake
	}{
		{
			name:     "empty",
			snakes:   []rules.Snake{},
			state:    map[string]SnakeState{},
			expected: []client.Snake{},
		},
		{
			name: "all properties",
			snakes: []rules.Snake{
				{ID: "one", Body: []rules.Point{{X: 3, Y: 3}, {X: 2, Y: 3}}, Health: 100},
			},
			state: map[string]SnakeState{
				"one": {
					ID:        "one",
					Name:      "ONE",
					URL:       "http://example1.com",
					Head:      "a",
					Tail:      "b",
					Color:     "#012345",
					LastMove:  "up",
					Character: '+',
					Latency:   time.Millisecond * 42,
				},
			},
			expected: []client.Snake{
				{
					ID:      "one",
					Name:    "ONE",
					Latency: "42",
					Health:  100,
					Body:    []client.Coord{{X: 3, Y: 3}, {X: 2, Y: 3}},
					Head:    client.Coord{X: 3, Y: 3},
					Length:  2,
					Shout:   "",
					Customizations: client.Customizations{
						Color: "#012345",
						Head:  "a",
						Tail:  "b",
					},
				},
			},
		},
		{
			name: "some eliminated",
			snakes: []rules.Snake{
				{
					ID:               "one",
					EliminatedCause:  rules.EliminatedByCollision,
					EliminatedOnTurn: 1,
					Body:             []rules.Point{{X: 3, Y: 3}},
				},
				{ID: "two", Body: []rules.Point{{X: 4, Y: 3}}},
			},
			state: map[string]SnakeState{
				"one": {ID: "one"},
				"two": {ID: "two"},
			},
			expected: []client.Snake{
				{
					ID:      "two",
					Latency: "0",
					Body:    []client.Coord{{X: 4, Y: 3}},
					Head:    client.Coord{X: 4, Y: 3},
					Length:  1,
				},
			},
		},
		{
			name: "all eliminated",
			snakes: []rules.Snake{
				{
					ID:               "one",
					EliminatedCause:  rules.EliminatedByCollision,
					EliminatedOnTurn: 1,
					Body:             []rules.Point{{X: 3, Y: 3}},
				},
			},
			state: map[string]SnakeState{
				"one": {ID: "one"},
			},
			expected: []client.Snake{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := convertRulesSnakes(test.snakes, test.state)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestBuildFrameEvent(t *testing.T) {
	tests := []struct {
		name        string
		boardState  *rules.BoardState
		snakeStates map[string]SnakeState
		expected    board.GameEvent
	}{
		{
			name:        "empty",
			boardState:  rules.NewBoardState(11, 11),
			snakeStates: map[string]SnakeState{},
			expected: board.GameEvent{
				EventType: board.EVENT_TYPE_FRAME,
				Data: board.GameFrame{
					Turn:   0,
					Snakes: []board.Snake{},
					Food:   []rules.Point{},
				},
			},
		},
		{
			name: "snake fields",
			boardState: rules.NewBoardState(19, 25).
				WithTurn(99).
				WithFood([]rules.Point{{X: 9, Y: 4}}).
				WithSnakes([]rules.Snake{
					{
						ID: "1",
						Body: []rules.Point{
							{X: 9, Y: 4},
							{X: 8, Y: 4},
							{X: 7, Y: 4},
						},
						Health:           97,
						EliminatedCause:  rules.EliminatedBySelfCollision,
						EliminatedOnTurn: 45,
						EliminatedBy:     "1",
					},
				}),
			snakeStates: map[string]SnakeState{
				"1": {
					URL:        "http://example.com",
					Name:       "One",
					ID:         "1",
					LastMove:   "left",
					Color:      "#ff00ff",
					Head:       "silly",
					Tail:       "default",
					Author:     "AUTHOR",
					Version:    "1.5",
					Error:      nil,
					StatusCode: 200,
					Latency:    54 * time.Millisecond,
				},
			},
			expected: board.GameEvent{
				EventType: board.EVENT_TYPE_FRAME,

				Data: board.GameFrame{
					Turn: 99,
					Snakes: []board.Snake{
						{
							ID:     "1",
							Name:   "One",
							Body:   []rules.Point{{X: 9, Y: 4}, {X: 8, Y: 4}, {X: 7, Y: 4}},
							Health: 97,
							Death: &board.Death{
								Cause:        rules.EliminatedBySelfCollision,
								Turn:         45,
								EliminatedBy: "1",
							},
							Color:         "#ff00ff",
							HeadType:      "silly",
							TailType:      "default",
							Latency:       "54",
							Author:        "AUTHOR",
							StatusCode:    200,
							Error:         "",
							IsBot:         false,
							IsEnvironment: false,
						},
					},
					Food: []rules.Point{{X: 9, Y: 4}},
				},
			},
		},
		{
			name: "snake errors",
			boardState: rules.NewBoardState(19, 25).
				WithSnakes([]rules.Snake{
					{
						ID: "bad_status",
					},
					{
						ID: "connection_error",
					},
				}),
			snakeStates: map[string]SnakeState{
				"bad_status": {
					StatusCode: 504,
					Latency:    54 * time.Millisecond,
				},
				"connection_error": {
					Error:   fmt.Errorf("error connecting to host"),
					Latency: 0,
				},
			},
			expected: board.GameEvent{
				EventType: board.EVENT_TYPE_FRAME,

				Data: board.GameFrame{
					Snakes: []board.Snake{
						{
							ID:         "bad_status",
							Latency:    "54",
							StatusCode: 504,
							Error:      "7:Bad HTTP status code 504",
						},
						{
							ID:         "connection_error",
							Latency:    "1",
							StatusCode: 0,
							Error:      "0:Error communicating with server",
						},
					},
					Food: []rules.Point{},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gameState := GameState{
				snakeStates: test.snakeStates,
			}
			actual := gameState.buildFrameEvent(test.boardState)
			require.Equalf(t, test.expected, actual, "%#v", actual)
		})
	}
}

func TestGetMoveForSnake(t *testing.T) {
	s1 := rules.Snake{ID: "one", Body: []rules.Point{{X: 3, Y: 3}}}
	s2 := rules.Snake{ID: "two", Body: []rules.Point{{X: 4, Y: 3}}}
	boardState := rules.NewBoardState(11, 11).WithSnakes([]rules.Snake{s1, s2})

	tests := []struct {
		name            string
		boardState      *rules.BoardState
		snakeState      SnakeState
		responseErr     error
		responseCode    int
		responseBody    string
		responseLatency time.Duration

		expectedSnakeState SnakeState
	}{
		{
			name:       "invalid URL",
			boardState: boardState,
			snakeState: SnakeState{
				ID:       "one",
				URL:      "",
				LastMove: rules.MoveLeft,
			},
			expectedSnakeState: SnakeState{
				ID:       "one",
				URL:      "",
				LastMove: rules.MoveLeft,
				Error:    errors.New(`parse "": empty url`),
			},
		},
		{
			name:       "error response",
			boardState: boardState,
			snakeState: SnakeState{
				ID:       "one",
				URL:      "http://example.com",
				LastMove: rules.MoveLeft,
			},
			responseErr: errors.New("connection error"),
			expectedSnakeState: SnakeState{
				ID:       "one",
				URL:      "http://example.com",
				LastMove: rules.MoveLeft,
				Error:    errors.New("connection error"),
			},
		},
		{
			name:       "bad response body",
			boardState: boardState,
			snakeState: SnakeState{
				ID:       "one",
				URL:      "http://example.com",
				LastMove: rules.MoveLeft,
			},
			responseCode:    200,
			responseBody:    `right`,
			responseLatency: 54 * time.Millisecond,
			expectedSnakeState: SnakeState{
				ID:         "one",
				URL:        "http://example.com",
				LastMove:   rules.MoveLeft,
				Error:      errors.New("invalid character 'r' looking for beginning of value"),
				StatusCode: 200,
				Latency:    54 * time.Millisecond,
			},
		},
		{
			name:       "bad move value",
			boardState: boardState,
			snakeState: SnakeState{
				ID:       "one",
				URL:      "http://example.com",
				LastMove: rules.MoveLeft,
			},
			responseCode:    200,
			responseBody:    `{"move": "north"}`,
			responseLatency: 54 * time.Millisecond,
			expectedSnakeState: SnakeState{
				ID:         "one",
				URL:        "http://example.com",
				LastMove:   rules.MoveLeft,
				StatusCode: 200,
				Latency:    54 * time.Millisecond,
			},
		},
		{
			name:       "bad status code",
			boardState: boardState,
			snakeState: SnakeState{
				ID:       "one",
				URL:      "http://example.com",
				LastMove: rules.MoveLeft,
			},
			responseCode:    500,
			responseLatency: 54 * time.Millisecond,
			expectedSnakeState: SnakeState{
				ID:         "one",
				URL:        "http://example.com",
				LastMove:   rules.MoveLeft,
				StatusCode: 500,
				Latency:    54 * time.Millisecond,
			},
		},
		{
			name:       "successful move",
			boardState: boardState,
			snakeState: SnakeState{
				ID:  "one",
				URL: "http://example.com",
			},
			responseCode:    200,
			responseBody:    `{"move": "right"}`,
			responseLatency: 54 * time.Millisecond,
			expectedSnakeState: SnakeState{
				ID:         "one",
				URL:        "http://example.com",
				LastMove:   rules.MoveRight,
				StatusCode: 200,
				Latency:    54 * time.Millisecond,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gameState := buildDefaultGameState()
			err := gameState.Initialize()
			require.NoError(t, err)
			gameState.snakeStates = map[string]SnakeState{test.snakeState.ID: test.snakeState}
			gameState.httpClient = stubHTTPClient{test.responseErr, test.responseCode, func(_ string) string { return test.responseBody }, test.responseLatency}

			nextSnakeState := gameState.getSnakeUpdate(test.boardState, test.snakeState)
			if test.expectedSnakeState.Error != nil {
				require.EqualError(t, nextSnakeState.Error, test.expectedSnakeState.Error.Error())
			} else {
				require.NoError(t, nextSnakeState.Error)
			}
			nextSnakeState.Error = test.expectedSnakeState.Error
			require.Equal(t, test.expectedSnakeState, nextSnakeState)
		})
	}
}

func TestCreateNextBoardState(t *testing.T) {
	s1 := rules.Snake{ID: "one", Body: []rules.Point{{X: 3, Y: 3}}}
	boardState := rules.NewBoardState(11, 11).WithSnakes([]rules.Snake{s1})
	snakeState := SnakeState{
		ID:  s1.ID,
		URL: "http://example.com",
	}

	t.Run("concurrently", func(t *testing.T) {
		gameState := buildDefaultGameState()
		err := gameState.Initialize()
		require.NoError(t, err)
		gameState.snakeStates = map[string]SnakeState{s1.ID: snakeState}
		gameState.httpClient = stubHTTPClient{nil, 200, func(_ string) string { return `{"move": "right"}` }, 54 * time.Millisecond}

		gameOver, nextBoardState, err := gameState.createNextBoardState(boardState)
		require.NoError(t, err)
		require.False(t, gameOver)
		snakeState = gameState.snakeStates[s1.ID]

		require.NotNil(t, nextBoardState)
		require.Equal(t, nextBoardState.Turn, 1)
		require.Equal(t, nextBoardState.Snakes[0].Body[0], rules.Point{X: 4, Y: 3})
		require.Equal(t, snakeState.LastMove, rules.MoveRight)
		require.Equal(t, snakeState.StatusCode, 200)
		require.Equal(t, snakeState.Latency, 54*time.Millisecond)
	})

}

type StubRuleset struct {
	maxTurns int
	settings settings.Settings
}

func (ruleset StubRuleset) Name() string                { return "standard" }
func (ruleset StubRuleset) Settings() settings.Settings { return ruleset.settings }
func (ruleset StubRuleset) Execute(prevState *rules.BoardState, moves []rulesets.SnakeMove) (bool, *rules.BoardState, error) {
	return prevState.Turn >= ruleset.maxTurns, prevState, nil
}

type stubHTTPClient struct {
	err        error
	statusCode int
	body       func(url string) string
	latency    time.Duration
}

func (client stubHTTPClient) request(url string) (*http.Response, time.Duration, error) {
	if client.err != nil {
		return nil, client.latency, client.err
	}
	body := ioutil.NopCloser(bytes.NewBufferString(client.body(url)))

	response := &http.Response{
		Header:     make(http.Header),
		Body:       body,
		StatusCode: client.statusCode,
	}

	return response, client.latency, nil
}

func (client stubHTTPClient) Get(url string) (*http.Response, time.Duration, error) {
	return client.request(url)
}

func (client stubHTTPClient) Post(url string, contentType string, body io.Reader) (*http.Response, time.Duration, error) {
	return client.request(url)
}
