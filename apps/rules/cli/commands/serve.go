package commands

import (
	stdlog "log"
	"net/http"

	"github.com/spf13/cobra"
	log "github.com/spf13/jwalterweatherman"
)

func NewServeCommand() *cobra.Command {
	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Serve Battlesnake games.",
		Long:  "Serve Battlesnake games.",
		Run: func(cmd *cobra.Command, args []string) {
			mux := http.NewServeMux()
			mux.HandleFunc("/create", handleCreate)

			stdlog.Println("http server started on :5001")
			err := http.ListenAndServe(":5001", mux)
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
		gameID:        r.URL.Query().Get("gameID"),
	}

	gameState.Names = append(gameState.Names, "test-snake")
	gameState.URLs = append(gameState.URLs, "https://snapepy.onrender.com")

	if err := gameState.Initialize(); err != nil {
		log.ERROR.Fatalf("Error initializing game: %v", err)
	}

	defer w.Write([]byte(gameState.gameID))

	if err := gameState.Run(); err != nil {
		log.ERROR.Fatalf("Error running game: %v", err)
	}

}
