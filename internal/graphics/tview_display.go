package graphics

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TViewDisplay manages a flicker-free terminal UI using tview
type TViewDisplay struct {
	app         *tview.Application
	mainView    *tview.TextView
	running     bool
	mu          sync.Mutex
	stopChan    chan struct{}
}

// NewTViewDisplay creates a new tview-based display
func NewTViewDisplay() *TViewDisplay {
	return &TViewDisplay{
		stopChan: make(chan struct{}),
	}
}

// Start initializes and starts the tview application
func (td *TViewDisplay) Start() error {
	td.mu.Lock()
	defer td.mu.Unlock()

	if td.running {
		return fmt.Errorf("display already running")
	}

	// Create text view for main content
	td.mainView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(false).
		SetWordWrap(false).
		SetChangedFunc(func() {
			td.app.Draw()
		})

	// Create application
	td.app = tview.NewApplication().
		SetRoot(td.mainView, true).
		EnableMouse(false)

	td.running = true

	// Run app in goroutine
	go func() {
		if err := td.app.Run(); err != nil {
			panic(err)
		}
	}()

	// Give tview time to initialize
	time.Sleep(50 * time.Millisecond)

	return nil
}

// Stop gracefully stops the display
func (td *TViewDisplay) Stop() {
	td.mu.Lock()
	defer td.mu.Unlock()

	if !td.running {
		return
	}

	td.running = false
	close(td.stopChan)
	
	if td.app != nil {
		td.app.Stop()
	}
}

// IsRunning returns whether the display is running
func (td *TViewDisplay) IsRunning() bool {
	td.mu.Lock()
	defer td.mu.Unlock()
	return td.running
}

// UpdateRecording updates the display with recording information
func (td *TViewDisplay) UpdateRecording(
	sessionName string,
	telemetry *TelemetryDisplay,
	packetsRecorded uint64,
	bytesWritten uint64,
	errors uint64,
	elapsed time.Duration,
) {
	td.mu.Lock()
	defer td.mu.Unlock()

	if !td.running || td.mainView == nil {
		return
	}

	var content strings.Builder
	
	// Header
	content.WriteString("[yellow:b:]═══════════════════════════════════════════════════════[white]\n")
	content.WriteString("[yellow:b:]  🎮 RECORDING SESSION[white]\n")
	content.WriteString(fmt.Sprintf("[cyan]  📝 Session: %s[white]\n", sessionName))
	content.WriteString("[yellow:b:]═══════════════════════════════════════════════════════[white]\n\n")

	// Telemetry
	if telemetry != nil {
		content.WriteString("[cyan:b:][ LIVE TELEMETRY DATA 🏎️ ][white]\n\n")
		content.WriteString(td.formatTelemetry(telemetry))
	}

	// Stats
	content.WriteString("\n")
	content.WriteString(td.formatStats(packetsRecorded, bytesWritten, errors, elapsed, "recording"))
	
	// Controls
	content.WriteString("\n[yellow]💡 Press 'q' to stop recording[white]\n")

	td.mainView.SetText(content.String())
}

// UpdatePlayback updates the display with playback information
func (td *TViewDisplay) UpdatePlayback(
	filename string,
	speed float64,
	telemetry *TelemetryDisplay,
	packetsPlayed uint64,
	bytesSent uint64,
	elapsed time.Duration,
	isPaused bool,
) {
	td.mu.Lock()
	defer td.mu.Unlock()

	if !td.running || td.mainView == nil {
		return
	}

	var content strings.Builder
	
	// Header
	content.WriteString("[yellow:b:]═══════════════════════════════════════════════════════[white]\n")
	content.WriteString("[yellow:b:]  🎬 PLAYBACK SESSION[white]\n")
	content.WriteString(fmt.Sprintf("[cyan]  📁 File: %s[white]\n", filename))
	content.WriteString(fmt.Sprintf("[cyan]  ⚡ Speed: %.1fx[white]\n", speed))
	content.WriteString("[yellow:b:]═══════════════════════════════════════════════════════[white]\n\n")

	// Telemetry (always show)
	if telemetry != nil {
		if isPaused {
			content.WriteString("[red:b:]⏸  PAUSED - LAST TELEMETRY DATA[white]\n\n")
		} else {
			content.WriteString("[cyan:b:][ LIVE TELEMETRY DATA 🏎️ ][white]\n\n")
		}
		content.WriteString(td.formatTelemetry(telemetry))
	}

	// Stats (only show when not paused)
	if !isPaused {
		content.WriteString("\n")
		content.WriteString(td.formatStats(packetsPlayed, bytesSent, 0, elapsed, "playback"))
	}
	
	// Controls
	content.WriteString("\n[yellow]💡 Press 'q' to stop, 'p' to pause/resume[white]\n")

	td.mainView.SetText(content.String())
}

// formatTelemetry creates the telemetry display string
func (td *TViewDisplay) formatTelemetry(t *TelemetryDisplay) string {
	var content strings.Builder

	// Speed
	speedBar := td.createColorBar(int(t.Speed), 350, 20, "green")
	content.WriteString(fmt.Sprintf("🏁 Speed        │ %s  [white:b:]%.0f km/h[white]\n", speedBar, t.Speed))
	
	// Throttle
	throttleBar := td.createColorBar(int(t.Throttle*100), 100, 20, "green")
	content.WriteString(fmt.Sprintf("🚀 Throttle     │ %s  [green:b:]%.0f%%[white]\n", throttleBar, t.Throttle*100))
	
	// Brake
	brakeBar := td.createColorBar(int(t.Brake*100), 100, 20, "red")
	content.WriteString(fmt.Sprintf("🛑 Brake        │ %s  [red:b:]%.0f%%[white]\n", brakeBar, t.Brake*100))
	
	// Engine RPM
	rpmBar := td.createColorBar(int(t.EngineRPM), 15000, 20, "yellow")
	content.WriteString(fmt.Sprintf("🏁 Engine RPM   │ %s  [yellow:b:]%d RPM[white]\n", rpmBar, t.EngineRPM))
	
	// Engine Temp
	tempColor := "green"
	if t.EngineTemp > 110 {
		tempColor = "red"
	} else if t.EngineTemp > 100 {
		tempColor = "yellow"
	}
	tempBar := td.createColorBar(int(t.EngineTemp)-50, 70, 20, tempColor)
	content.WriteString(fmt.Sprintf("🌡  Engine Temp  │ %s  [%s:b:]%d°C[white]\n", tempBar, tempColor, t.EngineTemp))
	
	// Tyre Temp
	tyreBar := td.createColorBar(int(t.TyreTempAvg)-20, 100, 20, "magenta")
	content.WriteString(fmt.Sprintf("🛞  Tyre Temp    │ %s  [magenta:b:]%d°C[white]\n", tyreBar, t.TyreTempAvg))
	
	// Fuel
	fuelBar := td.createColorBar(int(t.FuelLevel), 110, 20, "yellow")
	content.WriteString(fmt.Sprintf("⛽ Fuel Level   │ %s  [yellow:b:]%.1f kg[white]\n", fuelBar, t.FuelLevel))
	
	// ERS Energy
	ersPercent := (t.ERSEnergy / 4000000) * 100
	if ersPercent > 100 {
		ersPercent = 100
	}
	ersBar := td.createColorBar(int(ersPercent), 100, 20, "cyan")
	content.WriteString(fmt.Sprintf("⚡ ERS Energy   │ %s  [cyan:b:]%.0f%%[white]\n", ersBar, ersPercent))
	
	// Gear
	gearDisplay := ""
	if t.Gear == -1 {
		gearDisplay = "Reverse"
	} else if t.Gear == 0 {
		gearDisplay = "Neutral"
	} else {
		gearDisplay = fmt.Sprintf("Gear %d", int(t.Gear))
	}
	content.WriteString(fmt.Sprintf("🏁 Gearbox      │ [cyan:b:]%s[white]\n", gearDisplay))
	
	// DRS
	drsStatus := "[red]CLOSED[white]"
	if t.DRS {
		drsStatus = "[green:b:]OPEN[white]"
	}
	content.WriteString(fmt.Sprintf("💨 DRS Status   │ %s\n", drsStatus))
	
	content.WriteString("────────────────────────────────────────────────────\n")

	return content.String()
}

// formatStats creates the stats display string
func (td *TViewDisplay) formatStats(packets, bytes, errors uint64, elapsed time.Duration, mode string) string {
	var content strings.Builder

	content.WriteString(fmt.Sprintf("[cyan]📊 Packets: [white:b:]%d[white]  ", packets))
	content.WriteString(fmt.Sprintf("[cyan]💾 Data: [white:b:]%s[white]  ", formatBytes(bytes)))
	
	if mode == "recording" && errors > 0 {
		content.WriteString(fmt.Sprintf("[red]❌ Errors: [white:b:]%d[white]  ", errors))
	}
	
	content.WriteString(fmt.Sprintf("[cyan]⏱️  Time: [white:b:]%s[white]\n", formatDurationTview(elapsed)))

	return content.String()
}

// createColorBar creates a colored progress bar for tview
func (td *TViewDisplay) createColorBar(value, max, width int, color string) string {
	if value > max {
		value = max
	}
	if value < 0 {
		value = 0
	}
	
	filled := int(float64(value) / float64(max) * float64(width))
	var bar strings.Builder
	
	// Build filled portion with simple ASCII
	if filled > 0 {
		bar.WriteString(fmt.Sprintf("[%s]", color))
		for i := 0; i < filled; i++ {
			bar.WriteString("#")
		}
		bar.WriteString("[-]")
	}
	
	// Build empty portion with simple ASCII
	if filled < width {
		bar.WriteString("[gray]")
		for i := filled; i < width; i++ {
			bar.WriteString("-")
		}
		bar.WriteString("[-]")
	}
	
	return bar.String()
}

// formatBytes formats byte count into human-readable form
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatDuration formats duration into human-readable form (local version for tview)
func formatDurationTview(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

// HandleInput processes keyboard input for the display
func (td *TViewDisplay) HandleInput(handler func(key tcell.Key, ch rune)) {
	if td.app != nil {
		td.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			handler(event.Key(), event.Rune())
			return event
		})
	}
}
