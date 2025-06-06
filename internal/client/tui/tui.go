package tui

import (
	"fmt"
	"log"
	"os"

	"spacecraftsim/internal/client/core"
	"spacecraftsim/internal/parser"

	"github.com/rivo/tview"
)

// UI represents the terminal user interface
type UI struct {
	app      *tview.Application
	conn     *core.Connection
	logger   *core.Logger
	controls []Control
	grid     *tview.Grid
	logView  *tview.TextView
}

// New creates a new UI instance
func New(conn *core.Connection, logger *core.Logger) *UI {
	ui := &UI{
		app:     tview.NewApplication(),
		conn:    conn,
		logger:  logger,
		grid:    tview.NewGrid(),
		logView: tview.NewTextView().SetDynamicColors(true),
	}

	// Load device configuration
	config, err := LoadConfig("devices.yaml")
	if err != nil {
		log.Printf("Warning: Failed to load device config: %v", err)
		os.Exit(1)
		config = &Config{} // Use empty config as fallback
	}

	// Create controls from config
	ui.createControls(config)

	// Set up grid layout
	ui.setupLayout()

	// Set up message handler
	ui.conn.SetMessageHandler(func(msg core.Message) {
		// Update control values based on received messages
		for _, control := range ui.controls {
			if control.GetID() == msg.ID && len(msg.Values) > 0 {
				ui.app.QueueUpdateDraw(func() {
					if strValue, ok := msg.Values[0].(string); ok {
						control.SetValue(strValue)
					}
				})
			}
		}
		ui.logger.Log(core.LevelInfo, fmt.Sprintf("%s = %v", msg.ID, msg.Values))
	})

	// Set up response handler
	ui.conn.SetResponseHandler(func(resp parser.ResponseMessage) {
		if resp.Type == "error" {
			ui.app.QueueUpdateDraw(func() {
				ui.logger.Log(core.LevelError, fmt.Sprintf("Error from server: %s", resp.Error))
			})
		}
	})

	return ui
}

// Run starts the UI
func (ui *UI) Run() error {
	return ui.app.SetRoot(ui.grid, true).Run()
}

// createControls creates UI controls from configuration
func (ui *UI) createControls(config *Config) {
	for _, dev := range config.Devices {
		// Create local copy of device ID to avoid closure issues
		deviceID := dev.ID
		var control Control

		switch dev.Type {
		case TypeCheckbox:
			control = NewCheckboxControl(deviceID, dev.Label, func(checked bool) {
				ui.handleControlChange(deviceID, fmt.Sprintf("%v", checked))
			})

		case TypeSelector:
			control = NewSelectorControl(deviceID, dev.Label, dev.Options, func(option string, index int) {
				ui.handleControlChange(deviceID, option)
			})

		case TypeInput:
			control = NewInputControl(deviceID, dev.Label, func(text string) {
				ui.handleControlChange(deviceID, text)
			})

		default:
			log.Printf("Warning: Unknown device type %s for device %s", dev.Type, deviceID)
			continue
		}

		ui.controls = append(ui.controls, control)
	}
}

// setupLayout sets up the UI layout
func (ui *UI) setupLayout() {
	// Create a form for controls
	form := tview.NewForm()
	for _, control := range ui.controls {
		form.AddFormItem(control.GetFormItem())
	}

	// Set up grid layout
	ui.grid.SetRows(0, 10)
	ui.grid.SetColumns(0)
	ui.grid.AddItem(form, 0, 0, 1, 1, 0, 0, true)
	ui.grid.AddItem(ui.logView, 1, 0, 1, 1, 0, 0, false)

	// Set up log view
	ui.logView.SetBorder(true).SetTitle("Log")
	ui.logger.SetOutput(ui.logView)
}

// handleControlChange handles control value changes
func (ui *UI) handleControlChange(id, value string) {
	msg := core.Message{
		ID:     id,
		Values: []interface{}{value},
	}
	if err := ui.conn.SendMessage(msg); err != nil {
		ui.logger.Log(core.LevelError, fmt.Sprintf("Failed to send message: %v", err))
	} else {
		ui.logger.Log(core.LevelInfo, fmt.Sprintf("Sent: %s = %s", id, value))
	}
}
