package commands

import (
	"fmt"
	stdlog "log"
	"os"

	"github.com/spf13/cobra"

	log "github.com/spf13/jwalterweatherman"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "battlesnake",
	Short: "Battlesnake Command-Line Interface",
	Long:  "Tools and utilities for Battlesnake games.",
}

func Execute() {
	rootCmd.AddCommand(NewPlayCommand())
	rootCmd.AddCommand(NewServeCommand())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable debug logging")

	log.SetStdoutOutput(os.Stderr)
	log.SetFlags(stdlog.Ltime | stdlog.Lmicroseconds)
	if verbose {
		log.SetStdoutThreshold(log.LevelDebug)
	} else {
		log.SetStdoutThreshold(log.LevelInfo)
	}
}
