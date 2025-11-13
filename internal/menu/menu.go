package menu

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/pefman/golang-telemetry-recorder/internal/config"
	"github.com/pefman/golang-telemetry-recorder/internal/graphics"
	"github.com/pefman/golang-telemetry-recorder/internal/playback"
	"github.com/pefman/golang-telemetry-recorder/internal/recorder"
	"github.com/pefman/golang-telemetry-recorder/internal/session"
	"github.com/pefman/golang-telemetry-recorder/internal/telemetry"
)

const configFile = "config.json"

var (
	cfg    *config.Config
	reader *bufio.Reader
)

// Run starts the interactive menu system
func Run() error {
	reader = bufio.NewReader(os.Stdin)

	// Load configuration
	var err error
	cfg, err = config.Load(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Warning: Invalid configuration: %v\n", err)
		fmt.Println("Using default values...")
		cfg = config.NewDefaultConfig()
	}

	// Main menu loop
	for {
		showMainMenu()
		choice := readInput("Enter your choice: ")

		switch choice {
		case "1":
			if err := recordSession(); err != nil {
				fmt.Printf("Recording error: %v\n", err)
				pressEnterToContinue()
			}
		case "2":
			if err := playbackSession(); err != nil {
				fmt.Printf("Playback error: %v\n", err)
				pressEnterToContinue()
			}
		case "3":
			listRecordings()
			pressEnterToContinue()
		case "4":
			configureSettings()
		case "5":
			showStatus()
			pressEnterToContinue()
		case "6":
			fmt.Println("\nThank you for using F1 Telemetry Recorder!")
			return nil
		default:
			fmt.Println("Invalid choice. Please try again.")
			time.Sleep(1 * time.Second)
		}
	}
}

// showMainMenu displays the main menu
func showMainMenu() {
	clearScreen()
	fmt.Println("==============================================")
	fmt.Println("  F1 TELEMETRY RECORDER - MAIN MENU")
	fmt.Println("==============================================")
	fmt.Println()
	fmt.Println("  1. Start Recording")
	fmt.Println("  2. Playback Recording")
	fmt.Println("  3. List Recordings")
	fmt.Println("  4. Configure Settings")
	fmt.Println("  5. View Status")
	fmt.Println("  6. Exit")
	fmt.Println()
}

// recordSession handles recording telemetry data
func recordSession() error {
	clearScreen()
	fmt.Println("==============================================")
	fmt.Println("  RECORDING SESSION")
	fmt.Println("==============================================")
	fmt.Println()

	fmt.Println("🔍 Waiting for telemetry data to detect session info...")
	fmt.Println("   Start your F1 25 session now...")
	fmt.Println()

	// Create receiver first to detect session info
	receiverCfg := telemetry.ReceiverConfig{
		Port:       cfg.UDPPort,
		Address:    cfg.BindAddress,
		BufferSize: cfg.BufferSize,
		Timeout:    time.Duration(cfg.PacketTimeout) * time.Millisecond,
	}
	recv := telemetry.NewReceiver(receiverCfg)

	// Start receiver
	if err := recv.Start(); err != nil {
		return fmt.Errorf("failed to start receiver: %w", err)
	}
	defer recv.Stop()

	// Create a channel for session detection
	detectionChan := make(chan []byte, 100)
	sessionInfoChan := make(chan *session.SessionInfo, 1)
	
	// Start session info extraction goroutine
	go func() {
		info := session.ExtractSessionInfo(detectionChan)
		sessionInfoChan <- info
		close(detectionChan)
	}()

	// Buffer packets and send to detector
	var bufferedPackets []*telemetry.RecordedPacket
	timeout := time.After(30 * time.Second)
	var sessionInfo *session.SessionInfo

collectLoop:
	for {
		select {
		case packet, ok := <-recv.Packets():
			if !ok {
				break collectLoop
			}
			// Send copy to detector
			dataCopy := make([]byte, len(packet.Data))
			copy(dataCopy, packet.Data)
			select {
			case detectionChan <- dataCopy:
			default:
			}
			// Buffer the packet
			bufferedPackets = append(bufferedPackets, packet)
			
		case info := <-sessionInfoChan:
			sessionInfo = info
			break collectLoop
			
		case <-timeout:
			fmt.Println("\n⚠️  Timeout waiting for session data. Using default name.")
			sessionInfo = &session.SessionInfo{}
			break collectLoop
		}
	}

	// Generate session name
	sessionName := "session"
	if sessionInfo != nil && sessionInfo.HasInfo {
		sessionName = sessionInfo.GenerateFilename()
		fmt.Printf("\n✓ Session detected: %s\n", sessionInfo.String())
		fmt.Printf("✓ Recording as: %s\n", sessionName)
	} else {
		fmt.Println("\n⚠️  Could not detect session info, using default name")
		fmt.Print("   Enter custom name (or press Enter for 'session'): ")
		input := readInput("")
		if input != "" {
			sessionName = input
		}
	}

	// Create recorder with detected name
	rec, err := recorder.NewRecorder(cfg.RecordingDir, sessionName)
	if err != nil {
		return fmt.Errorf("failed to create recorder: %w", err)
	}

	// Start recorder
	if err := rec.Start(); err != nil {
		return fmt.Errorf("failed to start recorder: %w", err)
	}
	defer rec.Stop()

	// Record buffered packets first
	for _, packet := range bufferedPackets {
		if err := rec.RecordPacket(packet); err != nil {
			fmt.Printf("Error recording buffered packet: %v\n", err)
		}
	}

	// Initialize tview display
	display := graphics.NewTViewDisplay()
	if err := display.Start(); err != nil {
		return fmt.Errorf("failed to start display: %w", err)
	}
	defer display.Stop()

	// Start recording loop for new packets
	stopChan := make(chan struct{})
	userQuit := make(chan struct{})
	
	// Track latest telemetry data
	var latestTelemetry *telemetry.TelemetryData
	var playerCarIndex uint8
	
	// Get player car index from header
	if len(bufferedPackets) > 0 {
		playerCarIndex = bufferedPackets[0].Header.PlayerCarIndex
	}
	
	// Recording goroutine - must keep reading until recv stops
	go func() {
		for packet := range recv.Packets() {
			// Try to parse telemetry data
			if packet.Header.PacketID == 6 { // Car telemetry packet
				if td := telemetry.ParseCarTelemetryPacket(packet.Data, playerCarIndex); td != nil {
					latestTelemetry = telemetry.MergeTelemetryData(latestTelemetry, td)
				}
			} else if packet.Header.PacketID == 7 { // Car status packet
				if td := telemetry.ParseCarStatusPacket(packet.Data, playerCarIndex); td != nil {
					latestTelemetry = telemetry.MergeTelemetryData(latestTelemetry, td)
				}
			}
			
			if err := rec.RecordPacket(packet); err != nil {
				// Can't print errors in tview mode
			}
		}
	}()

	// Stats display goroutine with tview (flicker-free!)
	go func() {
		ticker := time.NewTicker(16 * time.Millisecond) // 60 FPS - ultra smooth!
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				// Update display if we have telemetry data
				if latestTelemetry == nil {
					continue
				}
				
				// Show live telemetry data
				telemetryDisplay := &graphics.TelemetryDisplay{
					Speed:       latestTelemetry.Speed,
					Throttle:    latestTelemetry.Throttle,
					Brake:       latestTelemetry.Brake,
					Gear:        latestTelemetry.Gear,
					EngineRPM:   latestTelemetry.EngineRPM,
					EngineTemp:  latestTelemetry.EngineTemp,
					TyreTempAvg: (latestTelemetry.TyreTemp[0] + latestTelemetry.TyreTemp[1] + 
								  latestTelemetry.TyreTemp[2] + latestTelemetry.TyreTemp[3]) / 4,
					FuelLevel:   latestTelemetry.FuelLevel,
					ERSEnergy:   latestTelemetry.ERSStoreEnergy,
					DRS:         latestTelemetry.DRS > 0,
				}
				
				// Show recording stats
				stats := rec.Stats()
				recvStats := recv.Stats()
				elapsed := time.Since(stats.StartTime)
				
				display.UpdateRecording(sessionName, telemetryDisplay, 
					stats.PacketsRecorded, stats.BytesWritten, 
					recvStats.Errors, elapsed)
			}
		}
	}()

	// Handle keyboard input
	display.HandleInput(func(key tcell.Key, ch rune) {
		if ch == 'q' || ch == 'Q' {
			select {
			case <-userQuit:
				// Already closing
			default:
				close(userQuit)
			}
		}
	})

	// Wait for user to stop
	<-userQuit
	close(stopChan)
	
	// Give time for goroutines to finish
	time.Sleep(100 * time.Millisecond)
	
	// Show completion message in terminal after tview stops
	clearScreen()
	stats := rec.Stats()
	duration := time.Since(stats.StartTime)
	graphics.ShowCompletionMessage("recording", stats.PacketsRecorded, 
		stats.BytesWritten, duration)
	
	fmt.Printf("\n💾 Output file: %s\n", rec.OutputPath())

	pressEnterToContinue()
	return nil
}

// playbackSession handles playback of recorded data
func playbackSession() error {
	clearScreen()
	fmt.Println("==============================================")
	fmt.Println("  PLAYBACK SESSION")
	fmt.Println("==============================================")
	fmt.Println()

	// List available recordings
	files, err := listRecordingFiles()
	if err != nil {
		return fmt.Errorf("failed to list recordings: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No recordings found.")
		pressEnterToContinue()
		return nil
	}

	// Display recordings
	for i, file := range files {
		info, _ := os.Stat(file)
		fmt.Printf("  %d. %s (%s, %s)\n",
			i+1,
			filepath.Base(file),
			formatFileSize(info.Size()),
			info.ModTime().Format("2006-01-02 15:04:05"))
	}
	fmt.Println()

	// Select recording
	choice := readInput("Select recording number: ")
	index, err := strconv.Atoi(choice)
	if err != nil || index < 1 || index > len(files) {
		return fmt.Errorf("invalid selection")
	}

	selectedFile := files[index-1]

	// Get playback settings
	fmt.Println()
	targetAddr := readInput(fmt.Sprintf("Target address (default: 127.0.0.1): "))
	if targetAddr == "" {
		targetAddr = "127.0.0.1"
	}

	targetPortStr := readInput(fmt.Sprintf("Target port (default: %d): ", cfg.UDPPort))
	targetPort := cfg.UDPPort
	if targetPortStr != "" {
		if p, err := strconv.Atoi(targetPortStr); err == nil {
			targetPort = p
		}
	}

	speedStr := readInput(fmt.Sprintf("Playback speed (default: %.1f): ", cfg.PlaybackSpeed))
	speed := cfg.PlaybackSpeed
	if speedStr != "" {
		if s, err := strconv.ParseFloat(speedStr, 64); err == nil && s > 0 {
			speed = s
		}
	}

	// Show playback initialization animation
	graphics.ShowWaveAnimation(1*time.Second, "🎬 Initializing Playback")

	// Create player
	player, err := playback.NewPlayer(selectedFile, targetAddr, targetPort, speed)
	if err != nil {
		return fmt.Errorf("failed to create player: %w", err)
	}

	// Start playback
	if err := player.Start(); err != nil {
		return fmt.Errorf("failed to start playback: %w", err)
	}
	defer player.Stop()

	// Initialize tview display
	display := graphics.NewTViewDisplay()
	if err := display.Start(); err != nil {
		return fmt.Errorf("failed to start display: %w", err)
	}
	defer display.Stop()

	// Stats display goroutine with tview (flicker-free!)
	stopChan := make(chan struct{})
	userQuit := make(chan struct{})
	
	// Track latest telemetry data
	var latestTelemetry *telemetry.TelemetryData
	var playerCarIndex uint8 = 0 // Will be updated from first packet
	
	// Process packets from playback for telemetry display
	go func() {
		for packet := range player.Packets() {
			// Update player car index from header
			if playerCarIndex == 0 {
				playerCarIndex = packet.Header.PlayerCarIndex
			}
			
			// Try to parse telemetry data
			if packet.Header.PacketID == 6 { // Car telemetry packet
				if td := telemetry.ParseCarTelemetryPacket(packet.Data, playerCarIndex); td != nil {
					latestTelemetry = telemetry.MergeTelemetryData(latestTelemetry, td)
				}
			} else if packet.Header.PacketID == 7 { // Car status packet
				if td := telemetry.ParseCarStatusPacket(packet.Data, playerCarIndex); td != nil {
					latestTelemetry = telemetry.MergeTelemetryData(latestTelemetry, td)
				}
			}
		}
	}()
	
	go func() {
		ticker := time.NewTicker(16 * time.Millisecond) // 60 FPS - ultra smooth!
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				// Show live telemetry data if available
				var telemetryDisplay *graphics.TelemetryDisplay
				if latestTelemetry != nil {
					telemetryDisplay = &graphics.TelemetryDisplay{
						Speed:       latestTelemetry.Speed,
						Throttle:    latestTelemetry.Throttle,
						Brake:       latestTelemetry.Brake,
						Gear:        latestTelemetry.Gear,
						EngineRPM:   latestTelemetry.EngineRPM,
						EngineTemp:  latestTelemetry.EngineTemp,
						TyreTempAvg: (latestTelemetry.TyreTemp[0] + latestTelemetry.TyreTemp[1] + 
									  latestTelemetry.TyreTemp[2] + latestTelemetry.TyreTemp[3]) / 4,
						FuelLevel:   latestTelemetry.FuelLevel,
						ERSEnergy:   latestTelemetry.ERSStoreEnergy,
						DRS:         latestTelemetry.DRS > 0,
					}
				}
				
				stats := player.Stats()
				elapsed := time.Since(stats.StartTime)
				
				display.UpdatePlayback(filepath.Base(selectedFile), speed, telemetryDisplay,
					stats.PacketsPlayed, stats.BytesSent, elapsed, player.IsPaused())
			}
		}
	}()

	// Handle keyboard input
	display.HandleInput(func(key tcell.Key, ch rune) {
		switch ch {
		case 'p', 'P':
			if player.IsPaused() {
				player.Resume()
			} else {
				player.Pause()
			}
		case 'q', 'Q':
			select {
			case <-userQuit:
				// Already closing
			default:
				close(userQuit)
				player.Stop()
			}
		}
	})

	// Wait for playback to finish or user to quit
	for player.IsRunning() {
		select {
		case <-userQuit:
			player.Stop()
		case <-time.After(100 * time.Millisecond):
			// Keep loop responsive
		}
	}

	close(stopChan)
	
	// Give time for goroutines to finish
	time.Sleep(100 * time.Millisecond)
	
	// Show completion message in terminal after tview stops
	clearScreen()
	stats := player.Stats()
	duration := time.Since(stats.StartTime)
	graphics.ShowCompletionMessage("playback", stats.PacketsPlayed, 
		stats.BytesSent, duration)

	pressEnterToContinue()
	return nil
}

// listRecordings displays all available recordings
func listRecordings() {
	clearScreen()
	fmt.Println("==============================================")
	fmt.Println("  AVAILABLE RECORDINGS")
	fmt.Println("==============================================")
	fmt.Println()

	files, err := listRecordingFiles()
	if err != nil {
		fmt.Printf("Error listing recordings: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No recordings found.")
		return
	}

	var totalSize int64
	for i, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		totalSize += info.Size()
		fmt.Printf("  %d. %s\n", i+1, filepath.Base(file))
		fmt.Printf("     Size: %s | Created: %s\n",
			formatFileSize(info.Size()),
			info.ModTime().Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	fmt.Printf("Total recordings: %d (%s)\n", len(files), formatFileSize(totalSize))
}

// configureSettings handles configuration menu
func configureSettings() {
	for {
		clearScreen()
		fmt.Println("==============================================")
		fmt.Println("  CONFIGURATION SETTINGS")
		fmt.Println("==============================================")
		fmt.Println()
		fmt.Printf("  1. UDP Port: %d\n", cfg.UDPPort)
		fmt.Printf("  2. Bind Address: %s\n", cfg.BindAddress)
		fmt.Printf("  3. Recording Directory: %s\n", cfg.RecordingDir)
		fmt.Printf("  4. Buffer Size: %d bytes\n", cfg.BufferSize)
		fmt.Printf("  5. Packet Timeout: %d ms\n", cfg.PacketTimeout)
		fmt.Printf("  6. Playback Speed: %.1fx\n", cfg.PlaybackSpeed)
		fmt.Println()
		fmt.Println("  7. Save Configuration")
		fmt.Println("  8. Reset to Defaults")
		fmt.Println("  9. Back to Main Menu")
		fmt.Println()

		choice := readInput("Enter your choice: ")

		switch choice {
		case "1":
			if val := readInput("Enter UDP port: "); val != "" {
				if port, err := strconv.Atoi(val); err == nil {
					cfg.UDPPort = port
				}
			}
		case "2":
			if val := readInput("Enter bind address: "); val != "" {
				cfg.BindAddress = val
			}
		case "3":
			if val := readInput("Enter recording directory: "); val != "" {
				cfg.RecordingDir = val
			}
		case "4":
			if val := readInput("Enter buffer size (bytes): "); val != "" {
				if size, err := strconv.Atoi(val); err == nil {
					cfg.BufferSize = size
				}
			}
		case "5":
			if val := readInput("Enter packet timeout (ms): "); val != "" {
				if timeout, err := strconv.Atoi(val); err == nil {
					cfg.PacketTimeout = timeout
				}
			}
		case "6":
			if val := readInput("Enter playback speed: "); val != "" {
				if speed, err := strconv.ParseFloat(val, 64); err == nil {
					cfg.PlaybackSpeed = speed
				}
			}
		case "7":
			if err := cfg.Save(configFile); err != nil {
				fmt.Printf("Error saving configuration: %v\n", err)
			} else {
				fmt.Println("✓ Configuration saved successfully!")
			}
			time.Sleep(1 * time.Second)
		case "8":
			cfg = config.NewDefaultConfig()
			fmt.Println("✓ Configuration reset to defaults!")
			time.Sleep(1 * time.Second)
		case "9":
			return
		}
	}
}

// showStatus displays current system status
func showStatus() {
	clearScreen()
	fmt.Println("==============================================")
	fmt.Println("  SYSTEM STATUS")
	fmt.Println("==============================================")
	fmt.Println()
	fmt.Printf("Configuration File: %s\n", configFile)
	fmt.Printf("UDP Port: %d\n", cfg.UDPPort)
	fmt.Printf("Bind Address: %s\n", cfg.BindAddress)
	fmt.Printf("Recording Directory: %s\n", cfg.RecordingDir)
	fmt.Println()

	// Check if recording directory exists
	if _, err := os.Stat(cfg.RecordingDir); os.IsNotExist(err) {
		fmt.Printf("⚠ Recording directory does not exist\n")
	} else {
		files, _ := listRecordingFiles()
		fmt.Printf("✓ Recording directory exists (%d recordings)\n", len(files))
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Printf("⚠ Configuration error: %v\n", err)
	} else {
		fmt.Println("✓ Configuration is valid")
	}
}

// Helper functions

func readInput(prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func clearScreen() {
	// Clear entire screen for menus
	fmt.Print("\033[H\033[2J")
}

func moveCursorHome() {
	// Move cursor to home position without clearing (for live updates)
	fmt.Print("\033[H")
}

func hideCursor() {
	fmt.Print("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}

func pressEnterToContinue() {
	fmt.Println()
	readInput("Press Enter to continue...")
}

func listRecordingFiles() ([]string, error) {
	pattern := filepath.Join(cfg.RecordingDir, "*.f1tr")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
