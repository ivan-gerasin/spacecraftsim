package tui

import (
	"fmt"
	"log"
	"os"

	"spacecraftsim/internal/client/core"

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

	return ui
}

// createControls creates UI controls from configuration
func (ui *UI) createControls(config *Config) {
	for _, dev := range config.Devices {
		// Create a new variable for each iteration to avoid closure issues
		device := dev
		var control Control

		switch device.Type {
		case TypeCheckbox:
			control = NewCheckboxControl(device.ID, device.Label, func(checked bool) {
				ui.handleControlChange(device.ID, fmt.Sprintf("%v", checked))
			})

		case TypeSelector:
			control = NewSelectorControl(device.ID, device.Label, device.Options, func(option string, index int) {
				ui.handleControlChange(device.ID, option)
			})

		case TypeInput:
			control = NewInputControl(device.ID, device.Label, func(text string) {
				ui.handleControlChange(device.ID, text)
			})

		default:
			log.Printf("Warning: Unknown device type %s for device %s", device.Type, device.ID)
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

// Run starts the UI
func (ui *UI) Run() error {
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

	// Start the application
	return ui.app.SetRoot(ui.grid, true).Run()
}
