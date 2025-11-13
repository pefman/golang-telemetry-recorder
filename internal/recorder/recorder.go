package recorder

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pefman/golang-telemetry-recorder/internal/telemetry"
)

// Recorder handles recording telemetry data to files
type Recorder struct {
	outputPath string
	file       *os.File
	mu         sync.Mutex
	stats      RecorderStats
	running    bool
}

// RecorderStats holds recording statistics
type RecorderStats struct {
	PacketsRecorded uint64
	BytesWritten    uint64
	StartTime       time.Time
	SessionName     string
}

// FileHeader is written at the start of recording files
type FileHeader struct {
	Magic       [4]byte   // "F1TR" (F1 Telemetry Recording)
	Version     uint16    // File format version
	Created     time.Time // Recording creation time
	Reserved    [32]byte  // Reserved for future use
}

// PacketEntry represents a packet in the recording file
type PacketEntry struct {
	Timestamp   int64  // Unix nanoseconds
	PacketSize  uint32 // Size of packet data
	PacketData  []byte // Raw packet data
}

// NewRecorder creates a new telemetry recorder
func NewRecorder(outputDir, sessionName string) (*Recorder, error) {
	// Create output directory if needed
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate output filename
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.f1tr", timestamp, sessionName)
	outputPath := filepath.Join(outputDir, filename)

	return &Recorder{
		outputPath: outputPath,
		stats: RecorderStats{
			SessionName: sessionName,
		},
	}, nil
}

// Start begins recording to file
func (r *Recorder) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		return fmt.Errorf("recorder already running")
	}

	// Open output file
	file, err := os.Create(r.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	r.file = file
	r.running = true
	r.stats.StartTime = time.Now()

	// Write file header
	if err := r.writeFileHeader(); err != nil {
		r.file.Close()
		return fmt.Errorf("failed to write file header: %w", err)
	}

	return nil
}

// Stop stops recording and closes the file
func (r *Recorder) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.running {
		return nil
	}

	r.running = false

	if r.file != nil {
		if err := r.file.Close(); err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}
	}

	return nil
}

// RecordPacket writes a packet to the recording file
func (r *Recorder) RecordPacket(packet *telemetry.RecordedPacket) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.running {
		return fmt.Errorf("recorder not running")
	}

	// Write timestamp (8 bytes)
	timestamp := packet.Timestamp.UnixNano()
	if err := binary.Write(r.file, binary.LittleEndian, timestamp); err != nil {
		return fmt.Errorf("failed to write timestamp: %w", err)
	}

	// Write packet size (4 bytes)
	packetSize := uint32(len(packet.Data))
	if err := binary.Write(r.file, binary.LittleEndian, packetSize); err != nil {
		return fmt.Errorf("failed to write packet size: %w", err)
	}

	// Write packet data
	n, err := r.file.Write(packet.Data)
	if err != nil {
		return fmt.Errorf("failed to write packet data: %w", err)
	}

	// Update stats
	r.stats.PacketsRecorded++
	r.stats.BytesWritten += uint64(n + 12) // data + timestamp + size

	return nil
}

// Stats returns current recorder statistics
func (r *Recorder) Stats() RecorderStats {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.stats
}

// IsRunning returns whether the recorder is active
func (r *Recorder) IsRunning() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.running
}

// OutputPath returns the output file path
func (r *Recorder) OutputPath() string {
	return r.outputPath
}

// writeFileHeader writes the file header
func (r *Recorder) writeFileHeader() error {
	header := FileHeader{
		Magic:   [4]byte{'F', '1', 'T', 'R'},
		Version: 1,
		Created: time.Now(),
	}

	// Write magic
	if _, err := r.file.Write(header.Magic[:]); err != nil {
		return err
	}

	// Write version
	if err := binary.Write(r.file, binary.LittleEndian, header.Version); err != nil {
		return err
	}

	// Write creation timestamp
	createdNano := header.Created.UnixNano()
	if err := binary.Write(r.file, binary.LittleEndian, createdNano); err != nil {
		return err
	}

	// Write reserved space
	if _, err := r.file.Write(header.Reserved[:]); err != nil {
		return err
	}

	return nil
}
