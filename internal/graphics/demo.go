package graphics

import (
	"fmt"
	"time"
)

// DemoRecording shows a demo of the recording interface
func DemoRecording() {
	ShowBootSequence()
	ShowRecordingHeader("demo_session")
	
	fmt.Println("📊 DEMO MODE - Simulating telemetry recording...\n")
	
	// Simulate recording for a few seconds
	for i := 0; i < 50; i++ {
		packets := uint64(i * 123)
		bytes := uint64(i * 1024 * 25)
		errors := uint64(0)
		elapsed := time.Duration(i*100) * time.Millisecond
		
		ShowLiveStats(packets, bytes, errors, elapsed, "recording")
		time.Sleep(100 * time.Millisecond)
	}
	
	fmt.Println("\n")
	ShowCompletionMessage("recording", 6150, 125952000, 5*time.Second)
	
	fmt.Println("\nPress Enter to continue...")
}

// DemoPlayback shows a demo of the playback interface
func DemoPlayback() {
	ShowWaveAnimation(1*time.Second, "🎬 Initializing Playback")
	ShowPlaybackHeader("monaco_quali_2025-11-13_19-45-30.f1tr", 1.5)
	
	fmt.Println("📊 DEMO MODE - Simulating telemetry playback...\n")
	
	// Simulate playback for a few seconds
	for i := 0; i < 50; i++ {
		packets := uint64(i * 98)
		bytes := uint64(i * 1024 * 20)
		errors := uint64(0)
		elapsed := time.Duration(i*100) * time.Millisecond
		
		// Simulate pause
		if i == 25 {
			ShowPausedState()
			time.Sleep(500 * time.Millisecond)
			fmt.Println("\n▶️  Resuming...")
			time.Sleep(200 * time.Millisecond)
		} else {
			ShowLiveStats(packets, bytes, errors, elapsed, "playback")
			time.Sleep(100 * time.Millisecond)
		}
	}
	
	fmt.Println("\n")
	ShowCompletionMessage("playback", 4900, 100352000, 5*time.Second)
	
	fmt.Println("\nPress Enter to continue...")
}

// ShowAllGraphics displays all available graphics
func ShowAllGraphics() {
	clearScreen()
	
	fmt.Println(Cyan + Bold + "╔════════════════════════════════════════════════════╗" + Reset)
	fmt.Println(Cyan + Bold + "║         F1 TELEMETRY RECORDER - GRAPHICS DEMO      ║" + Reset)
	fmt.Println(Cyan + Bold + "╚════════════════════════════════════════════════════╝" + Reset)
	fmt.Println()
	
	fmt.Println(Yellow + "1. Boot Sequence" + Reset)
	time.Sleep(500 * time.Millisecond)
	ShowBootSequence()
	
	fmt.Println("\nPress Enter for next demo...")
	fmt.Scanln()
	
	fmt.Println(Yellow + "\n2. Recording Demo" + Reset)
	time.Sleep(500 * time.Millisecond)
	DemoRecording()
	fmt.Scanln()
	
	fmt.Println(Yellow + "\n3. Playback Demo" + Reset)
	time.Sleep(500 * time.Millisecond)
	DemoPlayback()
	fmt.Scanln()
	
	fmt.Println(Yellow + "\n4. Car Animation" + Reset)
	time.Sleep(500 * time.Millisecond)
	fmt.Println("\n🏁 Race Start!")
	ShowCarAnimation(40)
	fmt.Println("🏁 Finish!")
	
	fmt.Println("\n" + Green + "✅ Graphics demo complete!" + Reset)
}
