package tui

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// DeviceType represents the type of UI control for a device
type DeviceType string

const (
	TypeCheckbox DeviceType = "checkbox"
	TypeSelector DeviceType = "selector"
	TypeInput    DeviceType = "input"
)

// DeviceConfig represents a device configuration from YAML
type DeviceConfig struct {
	ID      string     `yaml:"id"`
	Label   string     `yaml:"label"`
	Type    DeviceType `yaml:"type"`
	Options []string   `yaml:"options,omitempty"`
}

// Config represents the root configuration structure
type Config struct {
	Devices []DeviceConfig `yaml:"devices"`
}

// LoadConfig loads device configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate config
	for _, dev := range config.Devices {
		if dev.ID == "" {
			return nil, fmt.Errorf("device ID cannot be empty")
		}
		if dev.Label == "" {
			return nil, fmt.Errorf("device label cannot be empty")
		}
		if dev.Type == "" {
			return nil, fmt.Errorf("device type cannot be empty")
		}
		if dev.Type == TypeSelector && len(dev.Options) == 0 {
			return nil, fmt.Errorf("selector device must have options")
		}
	}

	return &config, nil
}
