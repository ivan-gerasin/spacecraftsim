package device

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// Sensor represents a read-only device that periodically generates values
type Sensor struct {
	*BaseDevice
	value     float64
	noise     float64
	lastValue float64
}

// NewSensor creates a new sensor device
func NewSensor(id string, initialValue, noise float64) *Sensor {
	s := &Sensor{
		BaseDevice: NewBaseDevice(id, time.Second),
		value:      initialValue,
		noise:      noise,
		lastValue:  initialValue,
	}
	s.AddTopic("sensors")
	return s
}

// HandleInput processes incoming messages
func (s *Sensor) HandleInput(msg Message) error {
	// Sensors are read-only, so they ignore input
	return nil
}

// Tick updates the sensor's value
func (s *Sensor) Tick() error {
	// Add some random noise to the value
	noise := (rand.Float64()*2 - 1) * s.noise
	s.value += noise

	// Only publish if the value has changed significantly
	if abs(s.value-s.lastValue) > s.noise/2 {
		s.lastValue = s.value
		msg := Message{
			ID:     s.id,
			Values: []interface{}{s.value},
			Time:   time.Now(),
			Source: s.id,
		}
		if err := s.bus.Publish("sensors", msg); err != nil {
			return fmt.Errorf("failed to publish sensor value: %w", err)
		}
		log.Printf("Sensor %s: %.2f", s.id, s.value)
	}

	return nil
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
