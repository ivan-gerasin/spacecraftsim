package ship

import (
	"fmt"
	"log"
	"sync"
	"time"

	"spacecraftsim/internal/bus"
	"spacecraftsim/internal/device"
)

// Ship represents the spacecraft system
type Ship struct {
	devices map[string]device.Device
	bus     *bus.MessageBus
	mu      sync.RWMutex
	stop    chan struct{}
}

// New creates a new ship system
func New() *Ship {
	return &Ship{
		devices: make(map[string]device.Device),
		bus:     bus.NewMessageBus(),
		stop:    make(chan struct{}),
	}
}

// RegisterDevice adds a device to the ship
func (s *Ship) RegisterDevice(dev device.Device) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.devices[dev.ID()]; exists {
		return fmt.Errorf("device with ID %s already exists", dev.ID())
	}

	if err := dev.Subscribe(s.bus); err != nil {
		return fmt.Errorf("failed to subscribe device %s: %w", dev.ID(), err)
	}

	s.devices[dev.ID()] = dev
	return nil
}

// Start begins the ship's operation
func (s *Ship) Start() {
	go s.run()
}

// Stop halts the ship's operation
func (s *Ship) Stop() {
	close(s.stop)
}

// HandleMessage processes an incoming message
func (s *Ship) HandleMessage(msg device.Message) error {
	s.mu.RLock()
	dev, exists := s.devices[msg.ID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("unknown device: %s", msg.ID)
	}

	return dev.HandleInput(msg)
}

// run executes the main ship loop
func (s *Ship) run() {
	// Use a faster base ticker for more precise timing
	baseTicker := time.NewTicker(10 * time.Millisecond)
	defer baseTicker.Stop()

	// Track last tick time for each device
	lastTicks := make(map[string]time.Time)

	for {
		select {
		case <-s.stop:
			return
		case <-baseTicker.C:
			s.mu.RLock()
			now := time.Now()
			for _, dev := range s.devices {
				tickRate := dev.GetTickRate()
				// Skip devices with zero tick rate
				if tickRate == 0 {
					continue
				}
				lastTick := lastTicks[dev.ID()]
				if now.Sub(lastTick) >= tickRate {
					if err := dev.Tick(); err != nil {
						log.Printf("Error ticking device %s: %v", dev.ID(), err)
					}
					lastTicks[dev.ID()] = now
				}
			}
			s.mu.RUnlock()
		}
	}
}
