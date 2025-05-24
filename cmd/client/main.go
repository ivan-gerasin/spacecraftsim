package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"spacecraftsim/internal/client/core"
	"spacecraftsim/internal/client/tui"
)

func main() {
	// Parse command line flags
	cliMode := flag.Bool("cli", false, "Use CLI mode instead of TUI")
	serverAddr := flag.String("server", "localhost:8080", "Server address")
	flag.Parse()

	// Create connection
	conn, err := core.NewConnection(*serverAddr)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Start heartbeat
	conn.StartHeartbeat(5 * time.Second)

	// Create logger
	logger := core.NewLogger(os.Stdout)

	if *cliMode {
		runCLI(conn, logger)
	} else {
		runTUI(conn, logger)
	}
}

func runCLI(conn *core.Connection, logger *core.Logger) {
	scanner := bufio.NewScanner(os.Stdin)
	logger.LogInfo("Connected to server. Enter messages (empty line to send):")

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		deviceID, values := core.ParseInput(line)
		if deviceID == "" {
			logger.LogError("Invalid input format. Expected: <device_id> <value1> <value2> ...")
			continue
		}

		if err := core.SendMessage(conn, deviceID, values); err != nil {
			logger.LogError("Failed to send message: %v", err)
		} else {
			logger.LogSuccess("Sent message to %s: %s", deviceID, strings.Join(values, " "))
		}
	}

	if err := scanner.Err(); err != nil {
		logger.LogError("Error reading input: %v", err)
	}
}

func runTUI(conn *core.Connection, logger *core.Logger) {
	ui := tui.New(conn, logger)
	if err := ui.Run(); err != nil {
		log.Fatalf("UI error: %v", err)
	}
}
