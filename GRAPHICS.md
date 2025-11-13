# Graphics & Visual Features

The F1 Telemetry Recorder now includes stunning animated graphics and visual feedback throughout the recording and playback experience!

## 🎨 Features

### 1. Boot Sequence Animation

When starting a recording session, you'll see an animated boot sequence:

```
╔════════════════════════════════════════════════════╗
║ 🏎️  F1 25 Telemetry Console — Booting Race Systems ║
╚════════════════════════════════════════════════════╝

[ POWER UNIT START SEQUENCE 🔧 ]
🏁  Engine RPM     │ ██████████▓▓▓▓▓▓▓▓▓▓  [10,420 RPM]
⚙️  Gearbox Sync    │ ██████████████▓▓▓▓▓▓  [Gear 5 Engaged]
🌡️  Engine Temp     │ ██████▓▓▓▓▓▓▓▓▓▓▓▓▓▓  [87°C Stable]

[ DRIVER INPUTS 🎮 ]
🚀  Throttle Pedal  │ ████████▓▓▓▓▓▓▓▓▓▓▓▓  [78%]
🛑  Brake Pressure   │ ███▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  [18%]
⚡  ERS Deployment   │ ████▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  [32%]
🔋  Battery Charge   │ ███████▓▓▓▓▓▓▓▓▓▓▓▓▓  [64%]

[ FUEL & TYRES 🛞 ]
⛽  Fuel Level       │ ████████▓▓▓▓▓▓▓▓▓▓▓▓  [68%]
🛞  Tyre Temp (Avg)  │ ██████▓▓▓▓▓▓▓▓▓▓▓▓▓▓  [91°C Grip Optimal]

[ CONNECTION 📡 ]
📡  Telemetry Link   │ ██████████▓▓▓▓▓▓▓▓▓▓  CONNECTED ✓

────────────────────────────────────────────────────
💬  Status: All systems nominal. Awaiting live packets.
────────────────────────────────────────────────────
```

### 2. Recording Interface

Live animated recording display with real-time statistics:

```
╔════════════════════════════════════════════════════╗
║ 🔴 RECORDING IN PROGRESS — DATA STREAM ACTIVE      ║
╚════════════════════════════════════════════════════╝

📝  Session: monaco_quali
⏺️   Press 'q' and Enter to stop recording...
────────────────────────────────────────────────────

🔴 ⠋ [ 00:05 ] Packets: 6,150 | Bytes: 120.21 MB | Errors: 0 | Time: 00:05 [▓▓▓▓▓░░░░░░░░░░]
```

Features:
- 🔴 **Live recording indicator**
- ⠋ **Spinning animation** (shows activity)
- 📊 **Real-time packet counter**
- 💾 **Live byte counter** with MB conversion
- ⚠️ **Error tracking**
- ⏱️ **Elapsed time**
- 📈 **Mini progress bar**

### 3. Playback Interface

Animated playback with controls:

```
╔════════════════════════════════════════════════════╗
║ ▶️  PLAYBACK MODE — REPLAYING TELEMETRY DATA       ║
╚════════════════════════════════════════════════════╝

📼  File: monaco_quali_2025-11-13_19-45-30.f1tr
⚡  Speed: 1.5x
🎮  Controls: 'p' = Pause/Resume | 'q' = Stop
────────────────────────────────────────────────────

▶️ ⠙ [ 00:03 ] Packets: 4,900 | Bytes: 95.66 MB | Errors: 0 | Time: 00:03 [▓▓▓▓▓▓▓░░░░░░░░]
```

When paused:
```
⏸️  PAUSED — Press 'p' to resume or 'q' to quit
```

### 4. Completion Screen

Celebrates successful operations with checkered flags:

```
╔════════════════════════════════════════════════════╗
║ ✅ RECORDING COMPLETED SUCCESSFULLY                 ║
╚════════════════════════════════════════════════════╝

📊  Total Packets: 6,150
💾  Total Bytes:   120.21 MB
⏱️   Duration:      00:05
📈  Avg Rate:      1230.0 packets/sec

  🏁 🏁 🏁 🏁 🏁 
```

### 5. Wave Animation

Smooth sine wave loading animation:

```
🎬 Initializing Playback ████▓▓▒▒░░░░▒▒▓▓████▓▓▒▒░░░░▒▒▓▓████▓▓▒▒░ ⠋
```

### 6. Car Animation

Racing car progress animation:

```
═════════════════════🏎️💨═══════════════════════════════
```

## 🎨 Color Scheme

- 🔴 **Red**: Recording, Critical states
- 🟢 **Green**: Success, Completed, Normal operation
- 🟡 **Yellow**: Warnings, Timeouts, Status
- 🔵 **Cyan**: Information, Headers
- 🟣 **Magenta**: Time, Duration
- ⚪ **White**: General text

## 📊 Live Statistics

All statistics update in real-time (10 times per second) for smooth animations:

- **Packets**: Total packets recorded/played
- **Bytes**: Data size with automatic MB conversion
- **Errors**: Error counter (for recording)
- **Time**: Elapsed time in HH:MM:SS or MM:SS format
- **Mini Bar**: Visual progress indicator

## 🎮 Interactive Elements

### Spinners
Multiple spinner styles for different contexts:
- ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏ (Braille patterns)

### Progress Bars
Different bar types:
- Full bars: ████████████
- Empty bars: ▓▓▓▓▓▓▓▓▓▓▓▓
- Gradient bars: █▓▒░

### Status Indicators
- 🔴 Recording
- ▶️ Playing
- ⏸️ Paused
- ✅ Completed
- ⚠️ Warning

## 🎪 Demo Mode

Try the new demo mode to see all graphics in action without needing F1 25:

1. Launch the application
2. Select **"6. Demo Graphics 🎨"**
3. Watch the animated demonstrations

The demo includes:
- Boot sequence animation
- Recording simulation
- Playback simulation with pause
- Car racing animation

## 💡 Technical Details

### Animation System
- **Update Rate**: 100ms (10 FPS)
- **Smooth Transitions**: Incremental value updates
- **Non-Blocking**: Goroutine-based rendering
- **Terminal Compatible**: ANSI color codes

### Performance
- Minimal CPU impact
- No impact on recording/playback performance
- Efficient string building
- Optimized rendering

### Compatibility
- ✅ Windows PowerShell
- ✅ Windows Terminal
- ✅ Command Prompt (limited colors)
- ✅ VS Code Terminal
- ⚠️ Legacy terminals (fallback to basic display)

## 🚀 Usage Tips

1. **Best Experience**: Use Windows Terminal for full color support
2. **Terminal Size**: Recommended minimum width: 60 characters
3. **Font**: Use a font with good emoji support (Cascadia Code, Consolas)
4. **Recording**: The animations don't affect data capture performance
5. **Playback**: Animations update independently of playback speed

## 🎉 Fun Facts

- The boot sequence shows simulated F1 telemetry values
- Progress bars animate smoothly during operations
- Spinners change every 100ms for visual feedback
- Checkered flags appear 5 times when completing operations
- Wave animations use sine wave mathematics
- Car animation smoothly crosses the screen

Enjoy the enhanced visual experience! 🏎️💨
