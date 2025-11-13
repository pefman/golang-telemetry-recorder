package telemetry

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Receiver handles UDP telemetry data reception
type Receiver struct {
	config   ReceiverConfig
	conn     *net.UDPConn
	packets  chan *RecordedPacket
	stopChan chan struct{}
	wg       sync.WaitGroup
	mu       sync.RWMutex
	running  bool
	stats    ReceiverStats
}

// ReceiverConfig holds receiver configuration
type ReceiverConfig struct {
	Port       int
	Address    string
	BufferSize int
	Timeout    time.Duration
}

// ReceiverStats holds receiver statistics
type ReceiverStats struct {
	PacketsReceived uint64
	BytesReceived   uint64
	Errors          uint64
	StartTime       time.Time
}

// NewReceiver creates a new telemetry receiver
func NewReceiver(config ReceiverConfig) *Receiver {
	return &Receiver{
		config:   config,
		packets:  make(chan *RecordedPacket, 100),
		stopChan: make(chan struct{}),
	}
}

// Start begins receiving telemetry packets
func (r *Receiver) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		return fmt.Errorf("receiver already running")
	}

	addr := &net.UDPAddr{
		IP:   net.ParseIP(r.config.Address),
		Port: r.config.Port,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to bind UDP socket: %w", err)
	}

	r.conn = conn
	r.running = true
	r.stats = ReceiverStats{StartTime: time.Now()}

	r.wg.Add(1)
	go r.receiveLoop()

	return nil
}

// Stop stops receiving telemetry packets
func (r *Receiver) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.running {
		return nil
	}

	close(r.stopChan)
	r.running = false

	if r.conn != nil {
		r.conn.Close()
	}

	r.wg.Wait()
	close(r.packets)

	return nil
}

// Packets returns the channel for received packets
func (r *Receiver) Packets() <-chan *RecordedPacket {
	return r.packets
}

// Stats returns current receiver statistics
func (r *Receiver) Stats() ReceiverStats {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.stats
}

// IsRunning returns whether the receiver is active
func (r *Receiver) IsRunning() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.running
}

// receiveLoop is the main receive loop
func (r *Receiver) receiveLoop() {
	defer r.wg.Done()

	buffer := make([]byte, r.config.BufferSize)

	for {
		select {
		case <-r.stopChan:
			return
		default:
			// Set read deadline for responsive shutdown
			r.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

			n, _, err := r.conn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue // Normal timeout, keep trying
				}
				if !r.IsRunning() {
					return // Connection closed during shutdown
				}
				r.incrementErrors()
				continue
			}

			if n > 0 {
				r.processPacket(buffer[:n])
			}
		}
	}
}

// processPacket processes a received packet
func (r *Receiver) processPacket(data []byte) {
	// Parse header
	header, err := ParseHeader(data)
	if err != nil {
		r.incrementErrors()
		return
	}

	// Create recorded packet
	packetData := make([]byte, len(data))
	copy(packetData, data)

	packet := &RecordedPacket{
		Timestamp: time.Now(),
		Data:      packetData,
		Header:    *header,
	}

	// Update stats
	r.mu.Lock()
	r.stats.PacketsReceived++
	r.stats.BytesReceived += uint64(len(data))
	r.mu.Unlock()

	// Send to channel (non-blocking)
	select {
	case r.packets <- packet:
	default:
		// Channel full, drop packet
		r.incrementErrors()
	}
}

// incrementErrors increments the error counter
func (r *Receiver) incrementErrors() {
	r.mu.Lock()
	r.stats.Errors++
	r.mu.Unlock()
}
