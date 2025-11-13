package graphics

import (
	"fmt"
	"strings"
)

// TelemetryDisplay holds current telemetry values for display
type TelemetryDisplay struct {
	Speed            float32
	Throttle         float32
	Brake            float32
	Gear             int8
	EngineRPM        uint16
	EngineTemp       uint16
	TyreTempAvg      uint8
	FuelLevel        float32
	ERSEnergy        float32
	DRS              bool
}

// ShowLiveTelemetry displays real-time telemetry data with bars
func ShowLiveTelemetry(td *TelemetryDisplay) {
	if td == nil {
		return
	}

	// Clear previous lines
	fmt.Print("\033[2K\r")
	
	fmt.Println(Cyan + "[ LIVE TELEMETRY DATA 🏎️ ]" + Reset)
	fmt.Println()
	
	// Speed
	speedBar := createBar(int(td.Speed), 350, 20, Green)
	fmt.Printf("🏁 Speed        │ %s  %s%.0f km/h%s\n", speedBar, Bold+White, td.Speed, Reset)
	
	// Throttle
	throttleBar := createBar(int(td.Throttle*100), 100, 20, Green)
	fmt.Printf("🚀 Throttle     │ %s  %s%.0f%%%s\n", throttleBar, Bold+Green, td.Throttle*100, Reset)
	
	// Brake
	brakeBar := createBar(int(td.Brake*100), 100, 20, Red)
	fmt.Printf("🛑 Brake        │ %s  %s%.0f%%%s\n", brakeBar, Bold+Red, td.Brake*100, Reset)
	
	// Engine RPM
	rpmBar := createBar(int(td.EngineRPM), 15000, 20, Yellow)
	fmt.Printf("🏁 Engine RPM   │ %s  %s%d RPM%s\n", rpmBar, Bold+Yellow, td.EngineRPM, Reset)
	
	// Engine Temp (adjusted scale: 50-120°C)
	tempColor := Green
	if td.EngineTemp > 110 {
		tempColor = Red
	} else if td.EngineTemp > 100 {
		tempColor = Yellow
	}
	tempBar := createBar(int(td.EngineTemp)-50, 70, 20, tempColor)
	fmt.Printf("🌡️  Engine Temp  │ %s  %s%d°C%s\n", tempBar, Bold+tempColor, td.EngineTemp, Reset)
	
	// Tyre Temp (adjusted scale: 20-120°C)
	tyreBar := createBar(int(td.TyreTempAvg)-20, 100, 20, Magenta)
	fmt.Printf("🛞  Tyre Temp    │ %s  %s%d°C%s\n", tyreBar, Bold+Magenta, td.TyreTempAvg, Reset)
	
	// Fuel
	fuelBar := createBar(int(td.FuelLevel), 110, 20, Yellow)
	fmt.Printf("⛽ Fuel Level   │ %s  %s%.1f kg%s\n", fuelBar, Bold+Yellow, td.FuelLevel, Reset)
	
	// ERS Energy (4,000,000 = 100%)
	ersPercent := (td.ERSEnergy / 4000000) * 100
	if ersPercent > 100 {
		ersPercent = 100
	}
	ersBar := createBar(int(ersPercent), 100, 20, Cyan)
	fmt.Printf("⚡ ERS Energy   │ %s  %s%.0f%%%s\n", ersBar, Bold+Cyan, ersPercent, Reset)
	
	// Gear
	gearDisplay := fmt.Sprintf("Gear %d", td.Gear)
	if td.Gear == -1 {
		gearDisplay = "Reverse"
	} else if td.Gear == 0 {
		gearDisplay = "Neutral"
	}
	fmt.Printf("⚙️  Gearbox      │ %s\n", Bold+Cyan+gearDisplay+Reset)
	
	// DRS
	drsStatus := Red + "CLOSED" + Reset
	if td.DRS {
		drsStatus = Green + Bold + "OPEN" + Reset
	}
	fmt.Printf("💨 DRS Status   │ %s\n", drsStatus)
	
	fmt.Println(strings.Repeat("─", 52))
}

// createBar creates a horizontal bar for values
func createBar(value, max, width int, color string) string {
	if value > max {
		value = max
	}
	if value < 0 {
		value = 0
	}
	
	filled := int(float64(value) / float64(max) * float64(width))
	bar := ""
	
	for i := 0; i < width; i++ {
		if i < filled {
			bar += color + "█" + Reset
		} else {
			bar += "▓"
		}
	}
	
	return bar
}
