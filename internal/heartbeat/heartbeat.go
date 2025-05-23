package heartbeat

import (
	"log"
	"net"
	"time"
)

const (
	heartbeatInterval    = 5 * time.Second
	maxReconnectAttempts = 3
	reconnectDelay       = 2 * time.Second
)

// HeartbeatMonitor manages connection health checks
type HeartbeatMonitor struct {
	conn            net.Conn
	stopChan        chan struct{}
	isConnected     bool
	reconnectCount  int
	onDisconnect    func()
	onReconnect     func()
	onMaxReconnects func()
	onConnUpdate    func(net.Conn)
}

// New creates a new heartbeat monitor
func New(conn net.Conn, onDisconnect, onReconnect, onMaxReconnects func(), onConnUpdate func(net.Conn)) *HeartbeatMonitor {
	return &HeartbeatMonitor{
		conn:            conn,
		stopChan:        make(chan struct{}),
		isConnected:     true,
		reconnectCount:  0,
		onDisconnect:    onDisconnect,
		onReconnect:     onReconnect,
		onMaxReconnects: onMaxReconnects,
		onConnUpdate:    onConnUpdate,
	}
}

// Start begins the heartbeat monitoring
func (h *HeartbeatMonitor) Start() {
	go h.monitor()
}

// Stop ends the heartbeat monitoring
func (h *HeartbeatMonitor) Stop() {
	close(h.stopChan)
}

// IsConnected returns the current connection status
func (h *HeartbeatMonitor) IsConnected() bool {
	return h.isConnected
}

// monitor runs the heartbeat loop
func (h *HeartbeatMonitor) monitor() {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-h.stopChan:
			return
		case <-ticker.C:
			if err := h.sendHeartbeat(); err != nil {
				h.handleDisconnect()
			}
		}
	}
}

// sendHeartbeat sends a heartbeat message
func (h *HeartbeatMonitor) sendHeartbeat() error {
	_, err := h.conn.Write([]byte("__heartbeat__\n"))
	return err
}

// handleDisconnect manages reconnection attempts
func (h *HeartbeatMonitor) handleDisconnect() {
	if !h.isConnected {
		return
	}

	h.isConnected = false
	h.onDisconnect()

	// Try to reconnect
	for h.reconnectCount < maxReconnectAttempts {
		log.Printf("Attempting to reconnect (attempt %d/%d)...", h.reconnectCount+1, maxReconnectAttempts)

		// Get the original address
		addr := h.conn.RemoteAddr().String()

		// Close the old connection
		h.conn.Close()

		// Wait before reconnecting
		time.Sleep(reconnectDelay)

		// Try to reconnect
		newConn, err := net.Dial("tcp", addr)
		if err == nil {
			h.conn = newConn
			h.isConnected = true
			h.reconnectCount = 0
			h.onReconnect()
			if h.onConnUpdate != nil {
				h.onConnUpdate(newConn)
			}
			return
		}

		h.reconnectCount++
	}

	// If we get here, we've failed to reconnect
	h.onMaxReconnects()
}
