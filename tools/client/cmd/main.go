package main

import (
	"flag"
	"log"
	"os"
	"time"

	"spacecraftsim/internal/client/core"
	"spacecraftsim/internal/client/tui"
	"spacecraftsim/tools/client"
)

func main() {
	// Parse command line flags
	cliMode := flag.Bool("cli", false, "Use CLI mode instead of TUI")
	serverAddr := flag.String("server", "localhost:8080", "Server address")
	flag.Parse()

	if *cliMode {
		// Use the existing CLI client
		cli, err := client.New(*serverAddr)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		defer cli.Close()
		cli.RunCLI()
	} else {
		// Create connection for TUI
		conn, err := core.NewConnection(*serverAddr)
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}
		defer conn.Close()

		// Start heartbeat
		conn.StartHeartbeat(5 * time.Second)

		// Create logger
		logger := core.NewLogger(os.Stdout)

		// Run TUI
		ui := tui.New(conn, logger)
		if err := ui.Run(); err != nil {
			log.Fatalf("UI error: %v", err)
		}
	}
}
