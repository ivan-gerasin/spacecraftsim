package tui

import (
	"fmt"
	"strings"

	"spacecraftsim/internal/client/core"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI represents the TUI application
type UI struct {
	app        *tview.Application
	deviceList *tview.List
	logView    *tview.TextView
	inputField *tview.InputField
	logger     *core.Logger
	conn       *core.Connection
}

// New creates a new TUI instance
func New(conn *core.Connection, logger *core.Logger) *UI {
	ui := &UI{
		app:    tview.NewApplication(),
		logger: logger,
		conn:   conn,
	}

	// Create the device list
	ui.deviceList = tview.NewList()
	ui.deviceList.ShowSecondaryText(false)
	ui.deviceList.SetBorder(true).
		SetTitle("Devices")

	// Create the log view
	ui.logView = tview.NewTextView()
	ui.logView.SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	ui.logView.SetBorder(true).
		SetTitle("Log")

	// Create the input field
	ui.inputField = tview.NewInputField()
	ui.inputField.SetLabel("Input: ")
	ui.inputField.SetBorder(true).
		SetTitle("Message")

	// Add some example devices (replace with actual device list)
	ui.deviceList.AddItem("logger1", "", 0, nil)
	ui.deviceList.AddItem("temp1", "", 0, nil)
	ui.deviceList.AddItem("press1", "", 0, nil)

	// Set up the layout
	grid := tview.NewGrid().
		SetRows(0, 3).
		SetColumns(30, 0).
		AddItem(ui.deviceList, 0, 0, 1, 1, 0, 0, true).
		AddItem(ui.logView, 0, 1, 1, 1, 0, 0, false).
		AddItem(ui.inputField, 1, 0, 1, 2, 0, 0, false)

	// Set up input handling
	ui.setupInputHandling()

	ui.app.SetRoot(grid, true)
	return ui
}

// setupInputHandling configures the input handling for the UI
func (ui *UI) setupInputHandling() {
	// Handle input field submission
	ui.inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			input := ui.inputField.GetText()
			if input == "" {
				return
			}

			// Get selected device
			_, deviceID := ui.deviceList.GetItemText(ui.deviceList.GetCurrentItem())
			values := strings.Fields(input)

			// Send the message
			if err := core.SendMessage(ui.conn, deviceID, values); err != nil {
				ui.logger.LogError("Failed to send message: %v", err)
			} else {
				ui.logger.LogSuccess("Sent message to %s: %s", deviceID, input)
			}

			// Clear the input field
			ui.inputField.SetText("")
		}
	})

	// Set up keyboard navigation
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			// Switch focus between components
			switch ui.app.GetFocus() {
			case ui.deviceList:
				ui.app.SetFocus(ui.inputField)
			case ui.inputField:
				ui.app.SetFocus(ui.deviceList)
			}
			return nil
		case tcell.KeyEsc:
			ui.app.Stop()
			return nil
		}
		return event
	})
}

// Run starts the TUI application
func (ui *UI) Run() error {
	return ui.app.Run()
}

// Log writes a message to the log view
func (ui *UI) Log(level core.LogLevel, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	ui.logView.Write([]byte(message + "\n"))
	ui.logView.ScrollToEnd()
}
