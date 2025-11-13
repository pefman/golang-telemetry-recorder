# Intelligent Session Detection

## Overview

The F1 Telemetry Recorder now features **automatic session detection** that intelligently extracts information from the F1 25 telemetry stream to create descriptive, meaningful filenames without manual input!

## 🧠 What Gets Detected

The system automatically extracts:

1. **Player Name** - Your driver name from the game
2. **Track Name** - Current circuit (Monaco, Silverstone, Spa, etc.)
3. **Session Type** - Practice, Qualifying, Race, Sprint, Time Trial, etc.
4. **Weather Conditions** - Clear, Light Rain, Heavy Rain, Storm, etc.

## 📝 Intelligent Filename Generation

Instead of generic names like `session_001.f1tr`, you now get descriptive filenames:

### Examples:
```
Monaco_Qualifying_Hamilton_Clear_2025-11-13_19-45-30.f1tr
Silverstone_Race_Verstappen_HeavyRain_2025-11-13_20-15-00.f1tr
Spa_Practice_YourName_LightCloud_2025-11-13_21-30-00.f1tr
Monza_TimeTrial_Player_Storm_2025-11-13_22-00-00.f1tr
```

### Filename Format:
```
{Track}_{SessionType}_{PlayerName}_{Weather}_{Timestamp}.f1tr
```

## 🔍 How It Works

### 1. Detection Process
When you start recording:
1. Receiver starts listening for F1 25 telemetry
2. System waits for **Session Packet** (Packet ID 1) - contains track, session type, weather
3. System waits for **Participants Packet** (Packet ID 4) - contains player name
4. Both packets are analyzed within 30 seconds
5. Filename is automatically generated

### 2. Packet Analysis

**Session Packet (ID 1):**
- Extracts track ID → Converts to track name
- Extracts session type → Converts to readable name  
- Extracts weather condition → Adds if not clear

**Participants Packet (ID 4):**
- Finds player's car index from packet header
- Extracts player name from participants array
- Sanitizes name for filesystem compatibility

### 3. Buffering
- All packets received during detection are buffered
- Once filename is determined, buffered packets are written first
- No telemetry data is lost during detection phase

## 📊 Supported Values

### Tracks (27 circuits)
- Melbourne, Shanghai, Bahrain, Catalunya
- **Monaco**, Montreal, **Silverstone**
- Hungaroring, **Spa**, **Monza**, Singapore
- Suzuka, Abu Dhabi, Texas, Brazil
- **Austria**, Mexico, **Baku**, Zandvoort
- **Imola**, Jeddah, **Miami**, Las Vegas, Losail
- Plus reverse layouts for Silverstone, Austria, Zandvoort

### Session Types
- **Practice** (P1, P2, P3, Short Practice)
- **Qualifying** (Q1, Q2, Q3, Short Qualifying, One-Shot)
- **Sprint Shootout** (SS1, SS2, SS3)
- **Race** (Race, Race 2, Race 3)
- **Time Trial**

### Weather Conditions
- Clear
- Light Cloud
- Overcast
- Light Rain
- Heavy Rain
- Storm

## 💡 Usage

### Automatic Detection (Recommended)
1. Select "Start Recording"
2. See message: "🔍 Waiting for telemetry data..."
3. Start or continue your F1 25 session
4. System displays: "✓ Session detected: Player: Hamilton | Track: Monaco | Session: Qualifying | Weather: Clear"
5. Recording begins with auto-generated filename

### Manual Override
If detection fails or times out (30 seconds):
1. System prompts: "Enter custom name (or press Enter for 'session')"
2. Type your own name or press Enter for default
3. Recording continues normally

## 🎯 Benefits

### Before (Manual Entry):
```
Enter session name: monaco_race
Recording as: 2025-11-13_19-45-30_monaco_race.f1tr
```

### After (Auto-Detection):
```
🔍 Waiting for telemetry data to detect session info...
✓ Session detected: Player: Hamilton | Track: Monaco | Session: Race | Weather: HeavyRain
✓ Recording as: Monaco_Race_Hamilton_HeavyRain
Recording as: Monaco_Race_Hamilton_HeavyRain_2025-11-13_19-45-30.f1tr
```

## ⚙️ Technical Details

### Packet Structure

**Session Packet Layout:**
```go
Offset 29:  uint8 weather         // 0=Clear, 1=LightCloud, etc.
Offset 30:  int8  trackTemp
Offset 31:  int8  airTemp
...
Offset 35:  uint8 sessionType     // 1=P1, 5=Q1, 15=Race, etc.
Offset 36:  int8  trackID         // 5=Monaco, 7=Silverstone, etc.
```

**Participants Packet Layout:**
```go
Header byte 27: playerCarIndex     // Index of player's car (0-21)
Offset 29:      numActiveCars
Offset 30+:     ParticipantData[22] array

Each ParticipantData (58 bytes):
  Offset 48: char name[32]         // UTF-8 null-terminated player name
```

### Name Sanitization
Player names are sanitized for filesystem compatibility:
- Spaces → Underscores
- Only alphanumeric, underscore, dash allowed
- UTF-8 characters are filtered
- Empty names default to "Player"

## 🚀 Performance

- **Detection Time**: < 1 second (typically)
- **Max Wait Time**: 30 seconds timeout
- **Packet Loss**: Zero - all packets buffered during detection
- **CPU Impact**: Minimal - only 2 packet types analyzed
- **Memory**: ~100 packets buffered (~100 KB)

## 🔮 Future Enhancements

Potential additions:
- [ ] Time of day extraction (day/night races)
- [ ] Lap number in filename
- [ ] Team name extraction
- [ ] Championship/career mode detection
- [ ] Multi-player session detection
- [ ] Custom filename templates
- [ ] Post-detection filename editing

## 📖 Example Workflow

### Scenario 1: Qualifying at Monaco
```
1. Start recorder → "Waiting for telemetry data..."
2. You're in F1 25 doing Monaco qualifying
3. System detects:
   - Track: Monaco
   - Session: Qualifying  
   - Player: YourName
   - Weather: Clear
4. Filename: Monaco_Qualifying_YourName_Clear_2025-11-13_19-45-30.f1tr
```

### Scenario 2: Rainy Race at Spa
```
1. Start recorder
2. Currently racing at Spa in heavy rain
3. System detects:
   - Track: Spa
   - Session: Race
   - Player: Hamilton
   - Weather: HeavyRain
4. Filename: Spa_Race_Hamilton_HeavyRain_2025-11-13_20-30-00.f1tr
```

### Scenario 3: Time Trial
```
1. Start recorder
2. Time trial mode at Silverstone
3. System detects:
   - Track: Silverstone
   - Session: TimeTrial
   - Player: Verstappen
   - Weather: Overcast
4. Filename: Silverstone_TimeTrial_Verstappen_Overcast_2025-11-13_21-00-00.f1tr
```

## ✅ Benefits Summary

✅ **No manual typing** - completely automatic
✅ **Descriptive filenames** - instantly know what's in the recording
✅ **Organized library** - easy to find specific sessions
✅ **Professional** - looks great in file listings
✅ **Sortable** - groups by track, session type, player
✅ **Zero data loss** - buffering ensures nothing is missed
✅ **Fast detection** - usually under 1 second
✅ **Fallback support** - manual override if needed

---

**Result**: Your recordings are now self-documenting with intelligent, descriptive names! 🎯
