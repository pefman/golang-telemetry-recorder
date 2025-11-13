package playback

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
	
	"github.com/pefman/golang-telemetry-recorder/internal/telemetry"
)

// Player handles playback of recorded telemetry data
type Player struct {
	filePath      string
	file          *os.File
	targetAddress string
	targetPort    int
	speed         float64
	conn          *net.UDPConn
	mu            sync.Mutex
	stats         PlayerStats
	running       bool
	paused        bool
	stopChan      chan struct{}
	packets       chan *telemetry.RecordedPacket
}

// PlayerStats holds playback statistics
type PlayerStats struct {
	PacketsPlayed uint64
	BytesSent     uint64
	StartTime     time.Time
	CurrentTime   time.Time
	RecordingTime time.Time
}

// NewPlayer creates a new telemetry player
func NewPlayer(filePath, targetAddress string, targetPort int, speed float64) (*Player, error) {
	return &Player{
		filePath:      filePath,
		targetAddress: targetAddress,
		targetPort:    targetPort,
		speed:         speed,
		stopChan:      make(chan struct{}),
		packets:       make(chan *telemetry.RecordedPacket, 100),
	}, nil
}

// Packets returns the channel for receiving parsed packets during playback
func (p *Player) Packets() <-chan *telemetry.RecordedPacket {
	return p.packets
}

// Start begins playback
func (p *Player) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return fmt.Errorf("player already running")
	}

	// Open recording file
	file, err := os.Open(p.filePath)
	if err != nil {
		return fmt.Errorf("failed to open recording file: %w", err)
	}
	p.file = file

	// Validate and skip file header
	if err := p.readFileHeader(); err != nil {
		p.file.Close()
		return fmt.Errorf("invalid recording file: %w", err)
	}

	// Setup UDP connection for sending
	addr := &net.UDPAddr{
		IP:   net.ParseIP(p.targetAddress),
		Port: p.targetPort,
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		p.file.Close()
		return fmt.Errorf("failed to create UDP connection: %w", err)
	}
	p.conn = conn

	p.running = true
	p.paused = false
	p.stats = PlayerStats{StartTime: time.Now()}

	// Start playback in goroutine
	go p.playbackLoop()

	return nil
}

// Stop stops playback
func (p *Player) Stop() error {
	p.mu.Lock()
	
	if !p.running {
		p.mu.Unlock()
		return nil
	}

	// Signal stop and mark as not running
	close(p.stopChan)
	p.running = false
	p.mu.Unlock()

	// Close connections (this will also wake up any blocking reads)
	if p.conn != nil {
		p.conn.Close()
	}

	if p.file != nil {
		p.file.Close()
	}
	
	// Wait a bit for the playback loop to exit
	time.Sleep(50 * time.Millisecond)
	
	// Now safe to close the packets channel
	close(p.packets)

	return nil
}

// Pause pauses playback
func (p *Player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paused = true
}

// Resume resumes playback
func (p *Player) Resume() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paused = false
}

// IsPaused returns whether playback is paused
func (p *Player) IsPaused() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.paused
}

// IsRunning returns whether playback is active
func (p *Player) IsRunning() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.running
}

// Stats returns current playback statistics
func (p *Player) Stats() PlayerStats {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.stats
}

// SetSpeed changes the playback speed
func (p *Player) SetSpeed(speed float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if speed > 0 {
		p.speed = speed
	}
}

// playbackLoop is the main playback loop
func (p *Player) playbackLoop() {
	var lastTimestamp int64 = 0

	for {
		select {
		case <-p.stopChan:
			return
		default:
			// Check if paused
			for p.IsPaused() {
				time.Sleep(100 * time.Millisecond)
				select {
				case <-p.stopChan:
					return
				default:
				}
			}

			// Read next packet
			timestamp, packetData, err := p.readPacket()
			if err != nil {
				if err == io.EOF {
					// Reached end of recording
					p.Stop()
					return
				}
				// Read error, stop playback
				p.Stop()
				return
			}

			// Calculate delay based on timestamp difference
			if lastTimestamp != 0 {
				timeDiff := time.Duration(timestamp - lastTimestamp)
				adjustedDelay := time.Duration(float64(timeDiff) / p.speed)

				if adjustedDelay > 0 {
					time.Sleep(adjustedDelay)
				}
			}

			// Send packet
			if err := p.sendPacket(packetData); err != nil {
				// Log error but continue
			}
			
			// Parse and send packet to channel for telemetry display (only if still running)
			if p.IsRunning() {
				if header, err := telemetry.ParseHeader(packetData); err == nil {
					packet := &telemetry.RecordedPacket{
						Timestamp: time.Unix(0, timestamp),
						Data:      packetData,
						Header:    *header,
					}
					select {
					case p.packets <- packet:
					default:
						// Channel full, skip (avoid blocking playback)
					}
				}
			}

			lastTimestamp = timestamp

			// Update stats
			p.mu.Lock()
			p.stats.PacketsPlayed++
			p.stats.BytesSent += uint64(len(packetData))
			p.stats.CurrentTime = time.Now()
			p.stats.RecordingTime = time.Unix(0, timestamp)
			p.mu.Unlock()
		}
	}
}

// readFileHeader reads and validates the file header
func (p *Player) readFileHeader() error {
	// Read magic
	magic := make([]byte, 4)
	if _, err := io.ReadFull(p.file, magic); err != nil {
		return err
	}
	if string(magic) != "F1TR" {
		return fmt.Errorf("invalid magic number")
	}

	// Read version
	var version uint16
	if err := binary.Read(p.file, binary.LittleEndian, &version); err != nil {
		return err
	}

	// Read creation timestamp
	var createdNano int64
	if err := binary.Read(p.file, binary.LittleEndian, &createdNano); err != nil {
		return err
	}

	// Skip reserved space
	reserved := make([]byte, 32)
	if _, err := io.ReadFull(p.file, reserved); err != nil {
		return err
	}

	return nil
}

// readPacket reads the next packet from the file
func (p *Player) readPacket() (int64, []byte, error) {
	// Read timestamp
	var timestamp int64
	if err := binary.Read(p.file, binary.LittleEndian, &timestamp); err != nil {
		return 0, nil, err
	}

	// Read packet size
	var packetSize uint32
	if err := binary.Read(p.file, binary.LittleEndian, &packetSize); err != nil {
		return 0, nil, err
	}

	// Read packet data
	packetData := make([]byte, packetSize)
	if _, err := io.ReadFull(p.file, packetData); err != nil {
		return 0, nil, err
	}

	return timestamp, packetData, nil
}

// sendPacket sends a packet via UDP
func (p *Player) sendPacket(data []byte) error {
	_, err := p.conn.Write(data)
	return err
}
