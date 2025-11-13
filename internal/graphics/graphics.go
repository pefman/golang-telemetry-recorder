package graphics

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// Colors for Windows console
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Bold    = "\033[1m"
)

// ShowBootSequence displays an animated boot sequence
func ShowBootSequence() {
	clearScreen()
	
	// Title
	fmt.Println(Cyan + Bold)
	fmt.Println("╔════════════════════════════════════════════════════╗")
	fmt.Println("║ 🏎️  F1 25 Telemetry Console — Booting Race Systems ║")
	fmt.Println("╚════════════════════════════════════════════════════╝")
	fmt.Println(Reset)
	time.Sleep(500 * time.Millisecond)

	// Power Unit Start
	fmt.Println(Yellow + "[ POWER UNIT START SEQUENCE 🔧 ]" + Reset)
	animateBar("🏁  Engine RPM    ", 10420, 15000, "RPM", Green)
	animateBar("⚙️  Gearbox Sync   ", 5, 8, "Gear 5 Engaged", Cyan)
	animateBar("🌡️  Engine Temp    ", 87, 120, "°C Stable", Yellow)
	fmt.Println()
	time.Sleep(300 * time.Millisecond)

	// Driver Inputs
	fmt.Println(Cyan + "[ DRIVER INPUTS 🎮 ]" + Reset)
	animateBar("🚀  Throttle Pedal ", 78, 100, "%", Green)
	animateBar("🛑  Brake Pressure  ", 18, 100, "%", Red)
	animateBar("⚡  ERS Deployment  ", 32, 100, "%", Magenta)
	animateBar("🔋  Battery Charge  ", 64, 100, "%", Yellow)
	fmt.Println()
	time.Sleep(300 * time.Millisecond)

	// Fuel & Tyres
	fmt.Println(Green + "[ FUEL & TYRES 🛞 ]" + Reset)
	animateBar("⛽  Fuel Level      ", 68, 100, "%", Yellow)
	animateBar("🛞  Tyre Temp (Avg) ", 91, 120, "°C Grip Optimal", Green)
	fmt.Println()
	time.Sleep(300 * time.Millisecond)

	// Connection
	fmt.Println(Magenta + "[ CONNECTION 📡 ]" + Reset)
	animateConnectionBar()
	fmt.Println()

	// Status
	fmt.Println(strings.Repeat("─", 52))
	fmt.Println(Green + "💬  Status: All systems nominal. Awaiting live packets." + Reset)
	fmt.Println(strings.Repeat("─", 52))
	fmt.Println()
	time.Sleep(800 * time.Millisecond)
}

// ShowRecordingHeader displays the recording session header
func ShowRecordingHeader(sessionName string) {
	clearScreen()
	fmt.Println(Red + Bold)
	fmt.Println("╔════════════════════════════════════════════════════╗")
	fmt.Println("║ 🔴 RECORDING IN PROGRESS — DATA STREAM ACTIVE      ║")
	fmt.Println("╚════════════════════════════════════════════════════╝")
	fmt.Println(Reset)
	fmt.Printf(Cyan+"📝  Session: %s\n"+Reset, sessionName)
	fmt.Println(Yellow + "⏺️   Press 'q' and Enter to stop recording..." + Reset)
	fmt.Println(strings.Repeat("─", 52))
	fmt.Println()
}

// ShowPlaybackHeader displays the playback session header
func ShowPlaybackHeader(filename string, speed float64) {
	clearScreen()
	fmt.Println(Green + Bold)
	fmt.Println("╔════════════════════════════════════════════════════╗")
	fmt.Println("║ ▶️  PLAYBACK MODE — REPLAYING TELEMETRY DATA        ║")
	fmt.Println("╚════════════════════════════════════════════════════╝")
	fmt.Println(Reset)
	fmt.Printf(Cyan+"📼  File: %s\n"+Reset, filename)
	fmt.Printf(Magenta+"⚡  Speed: %.1fx\n"+Reset, speed)
	fmt.Println(Yellow + "🎮  Controls: 'p' = Pause/Resume | 'q' = Stop" + Reset)
	fmt.Println(strings.Repeat("─", 52))
	fmt.Println()
}

// ShowLiveStats displays animated live statistics
func ShowLiveStats(packets, bytes, errors uint64, elapsed time.Duration, mode string) {
	icon := "🔴"
	color := Red
	if mode == "playback" {
		icon = "▶️"
		color = Green
	}

	// Create animated progress indicator
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	frame := frames[int(elapsed.Seconds())%len(frames)]

	// Format stats with bars
	packetBar := createMiniBar(int(packets%100), 100, 15)
	
	fmt.Printf("\r%s %s [ %s ] %s Packets: %s%d%s | Bytes: %s%.2f MB%s | Errors: %s%d%s | Time: %s%s%s   ",
		color+icon+Reset,
		Cyan+frame+Reset,
		formatDuration(elapsed),
		White,
		Green+Bold, packets, Reset,
		Cyan, float64(bytes)/(1024*1024), Reset,
		Yellow, errors, Reset,
		Magenta, formatDuration(elapsed), Reset+packetBar)
}

// ShowPausedState displays the paused state
func ShowPausedState() {
	fmt.Printf("\r%s ⏸️  PAUSED %s — Press 'p' to resume or 'q' to quit                                    ",
		Yellow+Bold, Reset)
}

// ShowCompletionMessage displays a completion message with animation
func ShowCompletionMessage(mode string, packets uint64, bytes uint64, duration time.Duration) {
	fmt.Println("\n")
	
	if mode == "recording" {
		fmt.Println(Green + Bold)
		fmt.Println("╔════════════════════════════════════════════════════╗")
		fmt.Println("║ ✅ RECORDING COMPLETED SUCCESSFULLY                 ║")
		fmt.Println("╚════════════════════════════════════════════════════╝")
		fmt.Println(Reset)
	} else {
		fmt.Println(Green + Bold)
		fmt.Println("╔════════════════════════════════════════════════════╗")
		fmt.Println("║ ✅ PLAYBACK COMPLETED                               ║")
		fmt.Println("╚════════════════════════════════════════════════════╝")
		fmt.Println(Reset)
	}

	fmt.Printf(Cyan+"📊  Total Packets: %s%d%s\n"+Reset, Bold, packets, Reset)
	fmt.Printf(Cyan+"💾  Total Bytes:   %s%.2f MB%s\n"+Reset, Bold, float64(bytes)/(1024*1024), Reset)
	fmt.Printf(Cyan+"⏱️   Duration:      %s%s%s\n"+Reset, Bold, formatDuration(duration), Reset)
	
	avgRate := float64(packets) / duration.Seconds()
	fmt.Printf(Cyan+"📈  Avg Rate:      %s%.1f packets/sec%s\n"+Reset, Bold, avgRate, Reset)
	
	fmt.Println()
	
	// Checkered flag animation
	time.Sleep(200 * time.Millisecond)
	fmt.Print("  ")
	for i := 0; i < 5; i++ {
		fmt.Print("🏁 ")
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println()
}

// ShowCarAnimation displays an animated racing car
func ShowCarAnimation(iterations int) {
	car := "🏎️💨"
	track := strings.Repeat("═", 50)
	
	for i := 0; i < iterations; i++ {
		pos := int(float64(i) / float64(iterations) * 48)
		fmt.Printf("\r%s%s%s", track[:pos], car, track[pos:])
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Println()
}

// Helper functions

func animateBar(label string, value, max int, unit string, color string) {
	fmt.Print(label + " │ ")
	
	steps := 20
	filled := int(float64(value) / float64(max) * float64(steps))
	
	for i := 0; i < steps; i++ {
		if i < filled {
			fmt.Print(color + "█" + Reset)
		} else {
			fmt.Print("▓")
		}
		time.Sleep(20 * time.Millisecond)
	}
	
	fmt.Printf("  [%s]\n", unit)
}

func animateConnectionBar() {
	fmt.Print("📡  Telemetry Link  │ ")
	
	steps := 20
	for i := 0; i < steps; i++ {
		if i < steps-1 {
			fmt.Print(Green + "█" + Reset)
		} else {
			fmt.Print(Green + "█" + Reset)
		}
		time.Sleep(30 * time.Millisecond)
	}
	
	fmt.Print("  " + Green + Bold + "CONNECTED ✓" + Reset + "\n")
}

func createMiniBar(value, max, width int) string {
	filled := int(float64(value) / float64(max) * float64(width))
	bar := " ["
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "▓"
		} else {
			bar += "░"
		}
	}
	bar += "]"
	return Cyan + bar + Reset
}

func formatDuration(d time.Duration) string {
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

func clearScreen() {
	// Clear entire screen
	fmt.Print("\033[H\033[2J")
}

// SpinnerFrames for loading animations
var SpinnerFrames = []string{
	"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏",
}

// ShowWaveAnimation shows a sine wave animation
func ShowWaveAnimation(duration time.Duration, label string) {
	start := time.Now()
	width := 50
	
	for time.Since(start) < duration {
		t := time.Since(start).Seconds()
		output := ""
		for i := 0; i < width; i++ {
			phase := float64(i)/float64(width)*2*math.Pi + t*2*math.Pi
			height := (math.Sin(phase) + 1) / 2
			
			if height > 0.7 {
				output += Cyan + "█" + Reset
			} else if height > 0.4 {
				output += Blue + "▓" + Reset
			} else if height > 0.2 {
				output += "▒"
			} else {
				output += "░"
			}
		}
		
		frameIndex := int(t*10) % len(SpinnerFrames)
		fmt.Printf("\r%s %s %s", Cyan+label+Reset, output, SpinnerFrames[frameIndex])
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Println()
}
