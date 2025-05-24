package device

import (
	"log"
)

// Logger represents a device that records incoming values
type Logger struct {
	*BaseDevice
	values []float64
}

// NewLogger creates a new logger device
func NewLogger(id string) *Logger {
	l := &Logger{
		BaseDevice: NewBaseDevice(id, 0),
		values:     make([]float64, 0),
	}
	l.AddTopic("logger")
	return l
}

// HandleInput processes incoming messages and logs the values
func (l *Logger) HandleInput(msg Message) error {
	for _, v := range msg.Values {
		if f, ok := v.(float64); ok {
			l.values = append(l.values, f)
			log.Printf("Logger %s received value: %.2f", l.id, f)
		}
	}
	return nil
}

// GetValues returns the recorded values
func (l *Logger) GetValues() []float64 {
	return l.values
}

func (l *Logger) Tick() error {
	return nil
}
