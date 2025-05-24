package bus

import (
	"fmt"
	"log"
	"sync"

	"spacecraftsim/internal/device"
)

// MessageBus implements the device.Bus interface
type MessageBus struct {
	subscribers map[string]map[device.Device]struct{}
	mu          sync.RWMutex
}

// NewMessageBus creates a new message bus
func NewMessageBus() *MessageBus {
	return &MessageBus{
		subscribers: make(map[string]map[device.Device]struct{}),
	}
}

// Publish sends a message to all subscribers of a topic
func (b *MessageBus) Publish(topic string, msg device.Message) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	subs, exists := b.subscribers[topic]
	if !exists {
		return nil // No subscribers is not an error
	}

	for dev := range subs {
		if err := dev.HandleInput(msg); err != nil {
			log.Printf("Error handling message for device %s: %v", dev.ID(), err)
		}
	}

	return nil
}

// Subscribe registers a device to receive messages on a topic
func (b *MessageBus) Subscribe(topic string, dev device.Device) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.subscribers[topic]; !exists {
		b.subscribers[topic] = make(map[device.Device]struct{})
	}

	b.subscribers[topic][dev] = struct{}{}
	return nil
}

// Unsubscribe removes a device's subscription to a topic
func (b *MessageBus) Unsubscribe(topic string, dev device.Device) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if subs, exists := b.subscribers[topic]; exists {
		delete(subs, dev)
		if len(subs) == 0 {
			delete(b.subscribers, topic)
		}
	}

	return nil
}

// Broadcast sends a message to all devices
func (b *MessageBus) Broadcast(msg device.Message) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for topic := range b.subscribers {
		if err := b.Publish(topic, msg); err != nil {
			return fmt.Errorf("error broadcasting to topic %s: %w", topic, err)
		}
	}

	return nil
}
