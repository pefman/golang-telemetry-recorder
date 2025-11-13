package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultUDPPort       = 20777 // F1 25 default UDP port
	DefaultBindAddress   = "0.0.0.0"
	DefaultRecordingDir  = "./recordings"
	DefaultBufferSize    = 65536
	DefaultPacketTimeout = 5000 // milliseconds
)

// Config holds the telemetry recorder configuration
type Config struct {
	// Network settings
	UDPPort     int    `json:"udp_port"`
	BindAddress string `json:"bind_address"`

	// Recording settings
	RecordingDir    string `json:"recording_dir"`
	AutoCreateDir   bool   `json:"auto_create_dir"`
	TimestampFormat string `json:"timestamp_format"`

	// Buffer settings
	BufferSize    int `json:"buffer_size"`
	PacketTimeout int `json:"packet_timeout"`

	// Playback settings
	PlaybackSpeed float64 `json:"playback_speed"` // 1.0 = real-time, 2.0 = 2x speed, etc.
}

// NewDefaultConfig returns a configuration with F1 25 defaults
func NewDefaultConfig() *Config {
	return &Config{
		UDPPort:         DefaultUDPPort,
		BindAddress:     DefaultBindAddress,
		RecordingDir:    DefaultRecordingDir,
		AutoCreateDir:   true,
		TimestampFormat: "2006-01-02_15-04-05",
		BufferSize:      DefaultBufferSize,
		PacketTimeout:   DefaultPacketTimeout,
		PlaybackSpeed:   1.0,
	}
}

// Load loads configuration from a JSON file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewDefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg := NewDefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

// Save saves configuration to a JSON file
func (c *Config) Save(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.UDPPort < 1 || c.UDPPort > 65535 {
		return fmt.Errorf("invalid UDP port: %d (must be 1-65535)", c.UDPPort)
	}

	if c.RecordingDir == "" {
		return fmt.Errorf("recording directory cannot be empty")
	}

	if c.BufferSize < 1024 {
		return fmt.Errorf("buffer size too small: %d (minimum 1024)", c.BufferSize)
	}

	if c.PlaybackSpeed <= 0 {
		return fmt.Errorf("invalid playback speed: %f (must be > 0)", c.PlaybackSpeed)
	}

	return nil
}
