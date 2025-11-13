# Quick Start Guide

## Installation & First Run

1. **Build the application** (if not already built):
   ```powershell
   go build -o f1-telemetry-recorder.exe
   ```

2. **Run the application**:
   ```powershell
   .\f1-telemetry-recorder.exe
   ```

## F1 25 Configuration

Before recording, configure F1 25:

1. Open F1 25
2. Navigate to: **Settings → Telemetry Settings**
3. Configure:
   - UDP Telemetry: **ON**
   - UDP Broadcast Mode: **ON**
   - UDP Port: **20777**
   - UDP Format: **2025**

## Recording Your First Session

1. Launch the F1 Telemetry Recorder
2. Select **"1. Start Recording"**
3. Enter a session name (e.g., "test_session")
4. Start your F1 25 game session
5. You'll see packets being received in real-time
6. Press **'q'** and Enter to stop recording
7. Your recording is saved in `./recordings/`

## Playing Back a Recording

1. Select **"2. Playback Recording"**
2. Choose a recording from the list
3. Accept defaults or customize:
   - Target Address: `127.0.0.1` (localhost)
   - Target Port: `20777` (F1 25 default)
   - Playback Speed: `1.0` (real-time)
4. Press **'q'** to stop playback

## Tips

- **First time**: Use default settings - they're optimized for F1 25
- **No packets?**: Check firewall settings and ensure F1 25 telemetry is enabled
- **Testing**: Record a short session first to verify everything works
- **Storage**: Each race generates approximately 50-200 MB of data
- **Playback**: You can play recordings back to localhost for analysis tools

## Common Issues

### "No packets received"
- Verify F1 25 is running with telemetry enabled
- Check Windows Firewall allows UDP on port 20777
- Ensure UDP port in config matches F1 25 settings

### "Failed to bind UDP socket"
- Port 20777 might be in use by another application
- Try changing the UDP port in configuration menu
- Update F1 25 to use the same port

### Recording directory errors
- Application will auto-create `./recordings` directory
- Ensure you have write permissions in the application folder

## Next Steps

- Explore the **Configuration Menu** (option 4) to customize settings
- Try different playback speeds for faster analysis
- Use recordings for developing analysis tools
- Share recordings with friends for comparison

## Need Help?

Check the full README.md for detailed documentation, or create an issue on GitHub.

Happy racing! 🏎️
