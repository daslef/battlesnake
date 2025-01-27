package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"rules"
	"rules/board"
	"rules/client"
	"rules/maps"
	"rules/rulesets"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	log "github.com/spf13/jwalterweatherman"
)

// Used to store state for each SnakeState while running a local game
type SnakeState struct {
	URL        string
	Name       string
	ID         string
	LastMove   string
	Character  rune
	Color      string
	Head       string
	Tail       string
	Author     string
	Version    string
	Error      error
	StatusCode int
	Latency    time.Duration
}

type GameState struct {
	Width           int
	Height          int
	Names           []string
	URLs            []string
	Timeout         int
	GameType        string
	Seed            int64
	ViewInBrowser   bool
	BoardURL        string
	FoodSpawnChance int

	// Internal game state
	settings    map[string]string
	snakeStates map[string]SnakeState
	gameID      string
	httpClient  TimedHttpClient
	ruleset     rulesets.Ruleset
	gameMap     maps.GameMap
	idGenerator func(int) string
}

func NewPlayCommand() *cobra.Command {
	gameState := &GameState{}

	var playCmd = &cobra.Command{
		Use:   "play",
		Short: "Play a game of Battlesnake locally.",
		Long:  "Play a game of Battlesnake locally.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gameState.Initialize(); err != nil {
				log.ERROR.Fatalf("Error initializing game: %v", err)
			}
			if err := gameState.Run(); err != nil {
				log.ERROR.Fatalf("Error running game: %v", err)
			}
		},
	}

	playCmd.Flags().IntVarP(&gameState.Width, "width", "W", 11, "Width of Board")
	playCmd.Flags().IntVarP(&gameState.Height, "height", "H", 11, "Height of Board")
	playCmd.Flags().StringArrayVarP(&gameState.Names, "name", "n", nil, "Name of Snake")
	playCmd.Flags().StringArrayVarP(&gameState.URLs, "url", "u", nil, "URL of Snake")

	playCmd.Flags().StringVarP(&gameState.GameType, "gametype", "g", "standard", "Type of Game Rules")
	playCmd.Flags().Int64VarP(&gameState.Seed, "seed", "r", time.Now().UTC().UnixNano(), "Random Seed")
	playCmd.Flags().BoolVar(&gameState.ViewInBrowser, "browser", false, "View the game in the browser using the Battlesnake game board")
	playCmd.Flags().StringVar(&gameState.BoardURL, "board-url", "https://board.battlesnake.com", "Base URL for the game board when using --browser")

	playCmd.Flags().IntVar(&gameState.FoodSpawnChance, "foodSpawnChance", 10, "Percentage chance of spawning a new food every round")

	playCmd.Flags().SortFlags = false

	return playCmd
}

// Setup a GameState once all the fields have been parsed from the command-line.
func (gameState *GameState) Initialize() error {
	// Generate game ID
	gameState.gameID = "tournament"

	// Set up HTTP client with request timeout
	gameState.Timeout = 500
	gameState.httpClient = timedHTTPClient{
		&http.Client{
			Timeout: time.Duration(gameState.Timeout) * time.Millisecond,
		},
	}

	gameState.gameMap = maps.StandardMap{}

	// Create settings object
	gameState.settings = map[string]string{
		rules.ParamFoodSpawnChance: fmt.Sprint(gameState.FoodSpawnChance),
	}

	// Build ruleset from settings
	ruleset := rulesets.NewRulesetBuilder().
		WithSeed(gameState.Seed).
		WithParams(gameState.settings).
		WithSolo(len(gameState.URLs) < 2).
		NamedRuleset(gameState.GameType)
	gameState.ruleset = ruleset

	// Initialize snake states as empty until we can ping the snake URLs
	gameState.snakeStates = map[string]SnakeState{}

	return nil
}

// Setup and run a full game.
func (gameState *GameState) Run() error {
	var gameOver bool
	var err error
	var winner SnakeState
	var isDraw = false

	// Setup local state for snakes
	gameState.snakeStates, err = gameState.buildSnakesFromOptions()
	if err != nil {
		return fmt.Errorf("error getting snake metadata: %w", err)
	}

	rand.Seed(gameState.Seed)

	gameOver, boardState, err := gameState.initializeBoardFromArgs()
	if err != nil {
		return fmt.Errorf("error initializing board: %w", err)
	}

	boardGame := board.Game{
		ID:     gameState.gameID,
		Status: "running",
		Width:  gameState.Width,
		Height: gameState.Height,
		Ruleset: map[string]string{
			rules.ParamGameType: gameState.GameType,
		},
		RulesetName: gameState.GameType,
		RulesStages: []string{},
		Map:         "standard",
	}

	boardServer := board.NewBoardServer(boardGame)

	if gameState.ViewInBrowser {
		serverURL, err := boardServer.Listen()
		if err != nil {
			return fmt.Errorf("error starting HTTP server: %w", err)
		}
		defer boardServer.Shutdown()
		log.INFO.Printf("Board server listening on %s", serverURL)

		boardURL := fmt.Sprintf(gameState.BoardURL+""+"?engine=%s&game=%s", serverURL, gameState.gameID)

		log.INFO.Printf("Opening board URL: %s", boardURL)

		// if err := browser.OpenURL(boardURL); err != nil {
		// 	log.ERROR.Printf("Failed to open browser: %v", err)
		// }

		// send turn zero to websocket server
		boardServer.SendEvent(gameState.buildFrameEvent(boardState))
	}

	gameState.printState(boardState)

	for !gameOver {
		gameOver, boardState, err = gameState.createNextBoardState(boardState)
		if err != nil {
			return fmt.Errorf("error processing game: %w", err)
		}

		if gameOver {
			break
		}

		// 	gameState.printState(boardState)

		if gameState.ViewInBrowser {
			boardServer.SendEvent(gameState.buildFrameEvent(boardState))
		}
	}

	if len(gameState.snakeStates) > 1 {
		// A draw is possible if there is more than one snake in the game.
		isDraw = true
	}

	for _, snake := range boardState.Snakes {
		snakeState := gameState.snakeStates[snake.ID]
		if snake.EliminatedCause == rules.NotEliminated {
			isDraw = false
			winner = snakeState
		}

		gameState.sendEndRequest(boardState, snakeState)
	}

	if isDraw {
		log.INFO.Printf("Game completed after %v turns. It was a draw.", boardState.Turn)
	} else if winner.Name != "" {
		log.INFO.Printf("Game completed after %v turns. %v was the winner.", boardState.Turn, winner.Name)
	} else {
		log.INFO.Printf("Game completed after %v turns.", boardState.Turn)
	}

	if gameState.ViewInBrowser {
		boardServer.SendEvent(board.GameEvent{
			EventType: board.EVENT_TYPE_GAME_END,
			Data:      boardGame,
		})
	}

	return nil
}

func (gameState *GameState) initializeBoardFromArgs() (bool, *rules.BoardState, error) {
	snakeIds := []string{}
	for _, snakeState := range gameState.snakeStates {
		snakeIds = append(snakeIds, snakeState.ID)
	}
	boardState, err := maps.SetupBoard(gameState.gameMap.ID(), gameState.ruleset.Settings(), gameState.Width, gameState.Height, snakeIds)
	if err != nil {
		return false, nil, fmt.Errorf("error initializing BoardState with map: %w", err)
	}
	gameOver, boardState, err := gameState.ruleset.Execute(boardState, nil)
	if err != nil {
		return false, nil, fmt.Errorf("error initializing BoardState with ruleset: %w", err)
	}

	for _, snakeState := range gameState.snakeStates {
		snakeRequest := gameState.getRequestBodyForSnake(boardState, snakeState)
		requestBody := serialiseSnakeRequest(snakeRequest)
		u, _ := url.ParseRequestURI(snakeState.URL)
		u.Path = path.Join(u.Path, "start")
		log.DEBUG.Printf("POST %s: %v", u, string(requestBody))
		_, _, err = gameState.httpClient.Post(u.String(), "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.WARN.Printf("Request to %v failed", u.String())
		}
	}
	return gameOver, boardState, nil
}

func (gameState *GameState) createNextBoardState(boardState *rules.BoardState) (bool, *rules.BoardState, error) {
	// apply PreUpdateBoard before making requests to snakes
	boardState, err := maps.PreUpdateBoard(gameState.gameMap, boardState, gameState.ruleset.Settings())
	if err != nil {
		return false, boardState, fmt.Errorf("error pre-updating board with game map: %w", err)
	}

	// get moves from snakes
	stateUpdates := make(chan SnakeState, len(gameState.snakeStates))

	var wg sync.WaitGroup
	for _, snakeState := range gameState.snakeStates {
		for _, snake := range boardState.Snakes {
			if snakeState.ID == snake.ID && snake.EliminatedCause == rules.NotEliminated {
				wg.Add(1)
				go func(snakeState SnakeState) {
					defer wg.Done()
					nextSnakeState := gameState.getSnakeUpdate(boardState, snakeState)
					stateUpdates <- nextSnakeState
				}(snakeState)
			}
		}
	}

	wg.Wait()
	close(stateUpdates)

	var moves []rulesets.SnakeMove
	for snakeState := range stateUpdates {
		gameState.snakeStates[snakeState.ID] = snakeState
		moves = append(moves, rulesets.SnakeMove{ID: snakeState.ID, Move: snakeState.LastMove})
	}

	gameOver, boardState, err := gameState.ruleset.Execute(boardState, moves)
	if err != nil {
		return false, boardState, fmt.Errorf("error updating board state from ruleset: %w", err)
	}

	// apply PostUpdateBoard after ruleset operates on snake moves
	boardState, err = maps.PostUpdateBoard(gameState.gameMap, boardState, gameState.ruleset.Settings())
	if err != nil {
		return false, boardState, fmt.Errorf("error post-updating board with game map: %w", err)
	}

	boardState.Turn += 1

	return gameOver, boardState, nil
}

func (gameState *GameState) getSnakeUpdate(boardState *rules.BoardState, snakeState SnakeState) SnakeState {
	snakeState.StatusCode = 0
	snakeState.Error = nil
	snakeState.Latency = 0

	snakeRequest := gameState.getRequestBodyForSnake(boardState, snakeState)
	requestBody := serialiseSnakeRequest(snakeRequest)

	u, err := url.ParseRequestURI(snakeState.URL)
	if err != nil {
		log.ERROR.Printf("Error parsing snake URL %#v: %v", snakeState.URL, err)
		snakeState.Error = err
		return snakeState
	}
	u.Path = path.Join(u.Path, "move")
	log.DEBUG.Printf("POST %s: %v", u, string(requestBody))
	res, responseTime, err := gameState.httpClient.Post(u.String(), "application/json", bytes.NewBuffer(requestBody))

	snakeState.Latency = responseTime

	if err != nil {
		log.WARN.Printf(
			"Request to %v failed\n"+
				"\tError: %s", u.String(), err)
		snakeState.Error = err
		return snakeState
	}

	snakeState.StatusCode = res.StatusCode

	if res.Body == nil {
		log.WARN.Printf(
			"Failed to parse response from %v\n"+
				"\tError: body is empty", u.String())
		return snakeState
	}
	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.WARN.Printf(
			"Failed to read response body from %v\n"+
				"\tError: %v", u.String(), readErr)
		snakeState.Error = readErr
		return snakeState
	}
	if res.StatusCode != http.StatusOK {
		log.WARN.Printf(
			"Got non-ok status code from %v\n"+
				"\tStatusCode: %d (expected %d)\n"+
				"\tBody: %q", u.String(), res.StatusCode, http.StatusOK, body)
		return snakeState
	}

	playerResponse := client.MoveResponse{}
	jsonErr := json.Unmarshal(body, &playerResponse)
	if jsonErr != nil {
		log.WARN.Printf(
			"Failed to decode JSON from %v\n"+
				"\tError: %v\n"+
				"\tBody: %q\n"+
				"\tSee https://docs.battlesnake.com/references/api#post-move", u.String(), jsonErr, body)
		snakeState.Error = jsonErr
		return snakeState
	}
	if playerResponse.Move != "up" && playerResponse.Move != "down" && playerResponse.Move != "left" && playerResponse.Move != "right" {
		log.WARN.Printf(
			"Failed to parse JSON data from %v\n"+
				"\tError: invalid move %q, valid moves are \"up\", \"down\", \"left\" or \"right\"\n"+
				"\tBody: %q\n"+
				"\tSee https://docs.battlesnake.com/references/api#post-move", u.String(), playerResponse.Move, body)
		return snakeState
	}

	snakeState.LastMove = playerResponse.Move

	return snakeState
}

func (gameState *GameState) sendEndRequest(boardState *rules.BoardState, snakeState SnakeState) {
	snakeRequest := gameState.getRequestBodyForSnake(boardState, snakeState)
	requestBody := serialiseSnakeRequest(snakeRequest)
	u, _ := url.ParseRequestURI(snakeState.URL)
	u.Path = path.Join(u.Path, "end")
	log.DEBUG.Printf("POST %s: %v", u, string(requestBody))
	_, _, err := gameState.httpClient.Post(u.String(), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.WARN.Printf("Request to %v failed", u.String())
	}
}

func (gameState *GameState) getRequestBodyForSnake(boardState *rules.BoardState, snakeState SnakeState) client.SnakeRequest {
	var youSnake rules.Snake
	for _, snk := range boardState.Snakes {
		if snakeState.ID == snk.ID {
			youSnake = snk
			break
		}
	}
	request := client.SnakeRequest{
		Game:  gameState.createClientGame(),
		Turn:  boardState.Turn,
		Board: convertStateToBoard(boardState, gameState.snakeStates),
		You:   convertRulesSnake(youSnake, snakeState),
	}
	return request
}

func (gameState *GameState) createClientGame() client.Game {
	return client.Game{
		ID:      gameState.gameID,
		Timeout: gameState.Timeout,
		Ruleset: client.Ruleset{
			Name:     gameState.ruleset.Name(),
			Settings: client.ConvertRulesetSettings(gameState.ruleset.Settings()),
		},
		Map: gameState.gameMap.ID(),
	}
}

func (gameState *GameState) buildSnakesFromOptions() (map[string]SnakeState, error) {
	bodyChars := []rune{'■', '⌀', '●', '☻', '◘', '☺', '□', '⍟'}
	var numSnakes int
	snakes := map[string]SnakeState{}
	numNames := len(gameState.Names)
	numURLs := len(gameState.URLs)
	if numNames > numURLs {
		numSnakes = numNames
	} else {
		numSnakes = numURLs
	}
	for i := int(0); i < numSnakes; i++ {
		var snakeName string
		var snakeURL string

		var id string
		if gameState.idGenerator != nil {
			id = gameState.idGenerator(i)
		} else {
			id = uuid.New().String()
		}

		if i < numNames {
			snakeName = gameState.Names[i]
		} else {
			log.DEBUG.Printf("name for URL %v is missing: a name will be generated automatically", gameState.URLs[i])
			snakeName = GenerateSnakeName()
		}

		if i < numURLs {
			u, err := url.ParseRequestURI(gameState.URLs[i])
			if err != nil {
				return nil, fmt.Errorf("url %v is not valid: %w", gameState.URLs[i], err)
			}
			snakeURL = u.String()
		} else {
			return nil, fmt.Errorf("url for name %v is missing", gameState.Names[i])
		}

		snakeState := SnakeState{
			Name: snakeName, URL: snakeURL, ID: id, LastMove: "up", Character: bodyChars[i%8],
		}
		var snakeErr error
		res, _, err := gameState.httpClient.Get(snakeURL)
		if err != nil {
			return nil, fmt.Errorf("snake metadata request to %v failed: %w", snakeURL, err)
		}

		snakeState.StatusCode = res.StatusCode

		if res.Body == nil {
			return nil, fmt.Errorf("empty response body from snake metadata URL: %v", snakeURL)
		}

		defer res.Body.Close()
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			return nil, fmt.Errorf("error reading from snake metadata URL %v: %w", snakeURL, readErr)
		}

		pingResponse := client.SnakeMetadataResponse{}
		jsonErr := json.Unmarshal(body, &pingResponse)
		if jsonErr != nil {
			return nil, fmt.Errorf("failed to parse response from %v: %w", snakeURL, jsonErr)
		}

		snakeState.Head = pingResponse.Head
		snakeState.Tail = pingResponse.Tail
		snakeState.Color = pingResponse.Color
		snakeState.Author = pingResponse.Author

		if snakeErr != nil {
			snakeState.Error = snakeErr
		}

		snakes[snakeState.ID] = snakeState

		log.INFO.Printf("Snake ID: %v URL: %v, Name: \"%v\"", snakeState.ID, snakeURL, snakeState.Name)
	}
	return snakes, nil
}

func (gameState *GameState) printState(boardState *rules.BoardState) {
	var aliveSnakeNames []string
	for _, snake := range boardState.Snakes {
		if snake.EliminatedCause == rules.NotEliminated {
			aliveSnakeNames = append(aliveSnakeNames, gameState.snakeStates[snake.ID].Name)
		}
	}
	log.INFO.Printf(
		"Turn: %d, Snakes Alive: [%v], Food: %d",
		boardState.Turn, strings.Join(aliveSnakeNames, ", "), len(boardState.Food),
	)
}

func (gameState *GameState) buildFrameEvent(boardState *rules.BoardState) board.GameEvent {
	snakes := []board.Snake{}

	for _, snake := range boardState.Snakes {
		snakeState := gameState.snakeStates[snake.ID]

		latencyMS := snakeState.Latency.Milliseconds()
		// round up latency of 0 to 1, to avoid legacy error display in board
		if latencyMS == 0 {
			latencyMS = 1
		}
		convertedSnake := board.Snake{
			ID:            snake.ID,
			Name:          snakeState.Name,
			Body:          snake.Body,
			Health:        snake.Health,
			Color:         snakeState.Color,
			HeadType:      snakeState.Head,
			TailType:      snakeState.Tail,
			Author:        snakeState.Author,
			StatusCode:    snakeState.StatusCode,
			IsBot:         false,
			IsEnvironment: false,
			Latency:       fmt.Sprint(latencyMS),
		}
		if snakeState.Error != nil {
			convertedSnake.Error = "0:Error communicating with server"
		} else if snakeState.StatusCode != http.StatusOK {
			convertedSnake.Error = fmt.Sprintf("7:Bad HTTP status code %d", snakeState.StatusCode)
		}
		if snake.EliminatedCause != rules.NotEliminated {
			convertedSnake.Death = &board.Death{
				Cause:        snake.EliminatedCause,
				Turn:         snake.EliminatedOnTurn,
				EliminatedBy: snake.EliminatedBy,
			}
		}
		snakes = append(snakes, convertedSnake)
	}

	gameFrame := board.GameFrame{
		Turn:   boardState.Turn,
		Snakes: snakes,
		Food:   boardState.Food,
	}

	return board.GameEvent{
		EventType: board.EVENT_TYPE_FRAME,
		Data:      gameFrame,
	}
}

func serialiseSnakeRequest(snakeRequest client.SnakeRequest) []byte {
	requestJSON, err := json.Marshal(snakeRequest)
	if err != nil {
		log.ERROR.Panicf("Error marshalling JSON from State: %v", err)
	}
	return requestJSON
}

func convertRulesSnake(snake rules.Snake, snakeState SnakeState) client.Snake {
	latencyMS := snakeState.Latency.Milliseconds()
	return client.Snake{
		ID:      snake.ID,
		Name:    snakeState.Name,
		Health:  snake.Health,
		Body:    client.CoordFromPointArray(snake.Body),
		Latency: fmt.Sprint(latencyMS),
		Head:    client.CoordFromPoint(snake.Body[0]),
		Length:  int(len(snake.Body)),
		Shout:   "",
		Customizations: client.Customizations{
			Head:  snakeState.Head,
			Tail:  snakeState.Tail,
			Color: snakeState.Color,
		},
	}
}

func convertRulesSnakes(snakes []rules.Snake, snakeStates map[string]SnakeState) []client.Snake {
	a := make([]client.Snake, 0)
	for _, snake := range snakes {
		if snake.EliminatedCause == rules.NotEliminated {
			a = append(a, convertRulesSnake(snake, snakeStates[snake.ID]))
		}
	}
	return a
}

func convertStateToBoard(boardState *rules.BoardState, snakeStates map[string]SnakeState) client.Board {
	return client.Board{
		Height: boardState.Height,
		Width:  boardState.Width,
		Food:   client.CoordFromPointArray(boardState.Food),
		Snakes: convertRulesSnakes(boardState.Snakes, snakeStates),
	}
}
