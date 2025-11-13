# F1 Telemetry Recorder

A powerful Golang-based telemetry recorder and playback tool for F1 25 (and other racing games that use UDP telemetry). Record your race sessions and play them back later for analysis, streaming, or testing.

## Features

- 🎮 **Real-time Recording**: Capture UDP telemetry data in real-time during gameplay
- 🧠 **Intelligent Session Detection**: Automatically extracts player name, track, session type, and weather
- 📹 **Playback**: Replay recorded sessions at various speeds
- ⚙️ **Fully Configurable**: Interactive menu system for easy configuration
- 🏎️ **F1 25 Optimized**: Pre-configured with F1 25 default settings
- 📊 **Live Telemetry Display**: Real-time telemetry visualization with speed, throttle, brake, RPM, temperatures, fuel, ERS, and DRS
- 🎯 **Flicker-Free Display**: Advanced tview/tcell rendering at 60 FPS with zero screen tearing
- 💾 **Efficient Storage**: Binary file format with timestamps for accurate playback
- 🎛️ **Playback Controls**: Pause, resume, and adjust playback speed on the fly with keyboard controls
- 🎨 **Animated Graphics**: Beautiful colored progress bars and live visual feedback
- 📝 **Smart Filenames**: Auto-generated descriptive filenames (Track_Session_Player_Weather)

## Quick Start

### Prerequisites

- Go 1.21 or higher
- F1 25 (or compatible racing game with UDP telemetry)

### Installation

1. Clone the repository:
```powershell
cd c:\Users\info\Documents\git\golang-telemetry-recorder
```

2. Build the application:
```powershell
go build -o f1-telemetry-recorder.exe
```

3. Run the application:
```powershell
.\f1-telemetry-recorder.exe
```

## F1 25 Setup

To enable telemetry output in F1 25:

1. Launch F1 25
2. Go to **Settings** → **Telemetry Settings**
3. Set **UDP Telemetry** to **On**
4. Set **UDP Broadcast Mode** to **On**
5. Set **UDP Port** to **20777** (default)
6. Set **UDP Format** to **2025** (latest format)

## Usage

### Main Menu

When you run the application, you'll see an interactive menu with the following options:

```
1. Start Recording      - Capture telemetry data from F1 25
2. Playback Recording   - Replay a previously recorded session
3. List Recordings      - View all saved recordings
4. Configure Settings   - Adjust application settings
5. View Status          - Check system status and configuration
6. Demo Graphics 🎨     - See animated graphics demonstration
7. Exit                 - Close the application
```

### Recording a Session

1. Select **"1. Start Recording"**
2. The recorder will automatically detect session information from F1 25:
   - **Player Name** - Your driver name
   - **Track** - Current circuit (e.g., Monaco, Silverstone)
   - **Session Type** - Practice, Qualifying, Race, etc.
   - **Weather** - Current weather conditions
3. Launch F1 25 and start your session (if not already started)
4. The system automatically generates an intelligent filename like:
   - `Monaco_Qualifying_YourName_LightRain_2025-11-13_19-45-30.f1tr`
   - `Silverstone_Race_Player_Clear_2025-11-13_20-15-00.f1tr`
5. **Watch live telemetry at 60 FPS** with smooth, flicker-free updates:
   - 🏁 Speed (km/h)
   - 🚀 Throttle (%)
   - 🛑 Brake (%)
   - 🏁 Engine RPM
   - 🌡️ Engine Temperature (°C)
   - 🛞 Tyre Temperature (°C)
   - ⛽ Fuel Level (kg)
   - ⚡ ERS Energy (%)
   - ⚙️ Gearbox (current gear)
   - 💨 DRS Status (OPEN/CLOSED)
6. See live statistics with packets, bytes, and elapsed time
7. Press **'q'** to stop recording
8. View completion summary with recording details

**Note**: If session info cannot be detected (offline menu, etc.), you can enter a custom name.

Recordings are saved in the `./recordings` directory with the format:
```
TrackName_SessionType_PlayerName_Weather_YYYY-MM-DD_HH-MM-SS.f1tr
```

### Playing Back a Recording

1. Select **"2. Playback Recording"**
2. Choose a recording from the list
3. Configure playback settings:
   - **Target Address**: Where to send UDP packets (default: 127.0.0.1)
   - **Target Port**: UDP port for playback (default: 20777)
   - **Playback Speed**: Speed multiplier (1.0 = real-time, 2.0 = 2x speed)
4. **Watch live telemetry during playback** with the same smooth 60 FPS display as recording
5. See live playback statistics showing packets sent, data volume, and elapsed time
6. Press **'p'** to pause/resume playback
7. Press **'q'** to stop playback
8. View completion summary with playback statistics

### Configuration Options

Access the configuration menu to customize:

- **UDP Port**: Port to listen on for telemetry (default: 20777)
- **Bind Address**: Network interface to bind to (default: 0.0.0.0)
- **Recording Directory**: Where to save recordings (default: ./recordings)
- **Buffer Size**: UDP receive buffer size (default: 65536 bytes)
- **Packet Timeout**: Timeout for packet reception (default: 5000 ms)
- **Playback Speed**: Default playback speed multiplier (default: 1.0)

Configuration is saved to `config.json` in the application directory.

## Graphics & Animations

The application features a modern terminal UI with flicker-free rendering powered by tview/tcell!

### Display Features

- � **60 FPS Updates**: Ultra-smooth telemetry display with zero screen flicker or tearing
- 📊 **Live Telemetry Bars**: Real-time colored progress bars for all telemetry metrics
- 🎨 **Color-Coded Values**: Green for throttle, red for brake, yellow for temperatures, cyan for ERS
- 💨 **Status Indicators**: DRS status (OPEN/CLOSED), current gear, fuel level, and more
- ⚙️ **Keyboard Controls**: Responsive 'q' to quit, 'p' to pause/resume during playback
- 🖥️ **PowerShell Compatible**: Optimized for Windows PowerShell and Windows Terminal

### Technical Details

Built with [tview](https://github.com/rivo/tview) and [tcell](https://github.com/gdamore/tcell) for professional terminal UI rendering:
- Double-buffered display prevents flickering
- Partial screen updates for optimal performance
- ANSI color support for rich visual feedback
- Works seamlessly in PowerShell, cmd, and modern terminals

Best experienced in Windows Terminal or any terminal with ANSI color support.

## Use Cases

### 📊 Data Analysis
Record race sessions and analyze telemetry data offline with external tools.

### 🎥 Content Creation
Replay sessions while streaming or recording gameplay for YouTube/Twitch content.

### 🧪 Testing & Development
Develop and test telemetry analysis tools without needing to play the game.

### 🏫 Coaching & Learning
Record and review sessions to improve driving techniques and setup changes.

### 🔄 Multi-Instance Playback
Send telemetry to multiple applications simultaneously for different analysis tools.

## File Format

Recordings use a custom binary format (`.f1tr`) with the following structure:

- **File Header**: Magic number, version, creation timestamp
- **Packet Entries**: Each entry contains:
  - Timestamp (int64, nanoseconds since epoch)
  - Packet size (uint32)
  - Raw packet data (variable length)

This format ensures accurate timing reproduction during playback.

## Troubleshooting

### No Packets Received

- Verify F1 25 telemetry is enabled in settings
- Check that UDP port matches F1 25 configuration (default: 20777)
- Ensure Windows Firewall allows UDP traffic on the configured port
- Try binding to a specific network interface instead of 0.0.0.0

### Playback Issues

- Verify the target application is listening on the specified port
- Check that the recording file is not corrupted
- Ensure sufficient system resources for the selected playback speed

### Configuration Errors

- Use "Reset to Defaults" option in the configuration menu
- Manually delete `config.json` to restore default settings
- Ensure recording directory has write permissions

## Advanced Usage

### Command Line (Future Enhancement)

The application currently uses an interactive menu. Future versions may support command-line arguments for automation:

```powershell
# Example future usage
.\f1-telemetry-recorder.exe record --name "my_session" --port 20777
.\f1-telemetry-recorder.exe playback --file "2025-11-13_recording.f1tr" --speed 2.0
```

### Network Recording

Record from another machine on your network by configuring the bind address:

1. On the recording machine, set bind address to its local IP (e.g., 192.168.1.100)
2. On F1 25 machine, set telemetry to broadcast to the recording machine's IP
3. Ensure network allows UDP traffic between machines

## Project Structure

```
golang-telemetry-recorder/
├── main.go                      # Application entry point
├── go.mod                       # Go module definition
├── internal/
│   ├── config/                  # Configuration management
│   │   └── config.go
│   ├── telemetry/               # Telemetry packet handling
│   │   ├── packet.go            # Packet structures and parsing
│   │   ├── parser.go            # F1 25 packet parsing
│   │   ├── receiver.go          # UDP receiver
│   │   └── errors.go            # Error definitions
│   ├── recorder/                # Recording functionality
│   │   └── recorder.go
│   ├── playback/                # Playback functionality
│   │   └── player.go
│   ├── session/                 # Session detection and naming
│   │   └── session.go
│   ├── graphics/                # Display and UI
│   │   ├── graphics.go          # ANSI colors and helpers
│   │   ├── telemetry.go         # Telemetry display (legacy)
│   │   └── tview_display.go    # Flicker-free tview UI
│   └── menu/                    # Interactive menu system
│       └── menu.go
└── recordings/                  # Default recording directory
```

## Contributing

Contributions are welcome! Areas for improvement:

- Support for other racing games (iRacing, ACC, rFactor 2, etc.)
- Data export to CSV/JSON formats for external analysis
- Web-based viewer/dashboard for telemetry visualization
- Packet filtering and selective recording by packet type
- Compression for recording files to save disk space
- Command-line interface mode for automation
- Additional telemetry metrics (brake temps, tyre wear, damage, etc.)
- Lap time analysis and sector comparisons

## License

This project is open source and available under the MIT License.

## Acknowledgments

- Built for the F1 25 community
- Inspired by the need for better telemetry analysis tools
- Thanks to all racing game developers who provide UDP telemetry APIs

## Support

For issues, questions, or feature requests, please create an issue on GitHub.

---

**Happy Racing! 🏎️💨**
