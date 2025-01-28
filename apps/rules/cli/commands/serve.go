package commands

import (
	stdlog "log"
	"math/rand"
	"net/http"

	"github.com/spf13/cobra"
	log "github.com/spf13/jwalterweatherman"

	"fmt"
	"time"

	"rules"
	"rules/board"
	"rules/maps"
	"rules/rulesets"
)

func NewServeCommand() *cobra.Command {
	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Serve Battlesnake games.",
		Long:  "Serve Battlesnake games.",
		Run: func(cmd *cobra.Command, args []string) {
			mux := http.NewServeMux()
			mux.HandleFunc("/create", handleCreate)

			stdlog.Println("http server started on http://localhost:5000")
			err := http.ListenAndServe(":5000", mux)
			if err != nil {
				stdlog.Fatal("ListenAndServe: ", err)
			}
		},
	}

	return serveCmd
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	gameState := &GameState{
		Width:         11,
		Height:        11,
		GameType:      "standard",
		ViewInBrowser: true,
		BoardURL:      "http://localhost:5173",
		gameMap:       maps.StandardMap{},
		Seed:          time.Now().UTC().UnixNano(),
		Timeout:       500,

		FoodSpawnChance: 15,
		settings: map[string]string{
			rules.ParamFoodSpawnChance: fmt.Sprint(15),
		},

		httpClient: timedHTTPClient{
			&http.Client{
				Timeout: time.Duration(500) * time.Millisecond,
			},
		},

		snakeStates: map[string]SnakeState{},
	}

	// r.URL.Query().Get("gameID")

	gameState.Names = append(gameState.Names, "test-snake")
	gameState.URLs = append(gameState.URLs, "https://snapepy.onrender.com")

	gameState.ruleset = rulesets.NewRulesetBuilder().
		WithSeed(gameState.Seed).
		WithParams(gameState.settings).
		WithSolo(len(gameState.URLs) < 2).
		NamedRuleset(gameState.GameType)

	if err := Run(gameState); err != nil {
		log.ERROR.Fatalf("Error running game: %v", err)
	}
}

func Run(gameState *GameState) error {
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

	serverURL, err := boardServer.Listen()
	if err != nil {
		return fmt.Errorf("error starting HTTP server: %w", err)
	}

	defer boardServer.Shutdown()
	
	log.INFO.Printf("Board server listening on %s", serverURL)

	// send turn zero to websocket server
	boardServer.SendEvent(gameState.buildFrameEvent(boardState))

	// gameState.printState(boardState)

	for !gameOver {
		gameOver, boardState, err = gameState.createNextBoardState(boardState)
		if err != nil {
			return fmt.Errorf("error processing game: %w", err)
		}

		if gameOver {
			break
		}

		// gameState.printState(boardState)
		boardServer.SendEvent(gameState.buildFrameEvent(boardState))
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

	boardServer.SendEvent(board.GameEvent{
		EventType: board.EVENT_TYPE_GAME_END,
		Data:      boardGame,
	})

	return nil
}
