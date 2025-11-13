package session

import (
	"fmt"
	"strings"
)

// SessionInfo holds extracted session information
type SessionInfo struct {
	PlayerName   string
	TrackName    string
	SessionType  string
	Weather      string
	TimeOfDay    string
	HasInfo      bool
}

// Track names mapped from track IDs
var trackNames = map[int8]string{
	0:  "Melbourne",
	2:  "Shanghai",
	3:  "Bahrain",
	4:  "Catalunya",
	5:  "Monaco",
	6:  "Montreal",
	7:  "Silverstone",
	9:  "Hungaroring",
	10: "Spa",
	11: "Monza",
	12: "Singapore",
	13: "Suzuka",
	14: "AbuDhabi",
	15: "Texas",
	16: "Brazil",
	17: "Austria",
	19: "Mexico",
	20: "Baku",
	26: "Zandvoort",
	27: "Imola",
	29: "Jeddah",
	30: "Miami",
	31: "LasVegas",
	32: "Losail",
	39: "Silverstone_Rev",
	40: "Austria_Rev",
	41: "Zandvoort_Rev",
}

// Session type names
var sessionTypes = map[uint8]string{
	0:  "Unknown",
	1:  "P1",
	2:  "P2",
	3:  "P3",
	4:  "Practice",
	5:  "Q1",
	6:  "Q2",
	7:  "Q3",
	8:  "Qualifying",
	9:  "OneShotQ",
	10: "SS1",
	11: "SS2",
	12: "SS3",
	13: "SprintShootout",
	14: "OneShotSS",
	15: "Race",
	16: "Race2",
	17: "Race3",
	18: "TimeTrial",
}

// Weather conditions
var weatherConditions = map[uint8]string{
	0: "Clear",
	1: "LightCloud",
	2: "Overcast",
	3: "LightRain",
	4: "HeavyRain",
	5: "Storm",
}

// ExtractSessionInfo extracts session information from packets
func ExtractSessionInfo(packets <-chan []byte) *SessionInfo {
	info := &SessionInfo{}
	
	sessionReceived := false
	participantsReceived := false
	
	// Process packets until we have both session and participant data
	for packetData := range packets {
		if len(packetData) < 29 {
			continue
		}
		
		packetID := packetData[6]
		
		switch packetID {
		case 1: // Session packet
			if !sessionReceived {
				parseSessionPacket(packetData, info)
				sessionReceived = true
			}
		case 4: // Participants packet
			if !participantsReceived {
				parseParticipantsPacket(packetData, info)
				participantsReceived = true
			}
		}
		
		// Exit early if we have all the info we need
		if sessionReceived && participantsReceived {
			info.HasInfo = true
			return info
		}
	}
	
	return info
}

// parseSessionPacket extracts info from session packet
func parseSessionPacket(data []byte, info *SessionInfo) {
	if len(data) < 100 {
		return
	}
	
	// Skip header (29 bytes)
	offset := 29
	
	weather := data[offset]
	info.Weather = weatherConditions[weather]
	if info.Weather == "" {
		info.Weather = "Unknown"
	}
	
	offset += 1 // m_weather
	offset += 1 // m_trackTemperature
	offset += 1 // m_airTemperature
	offset += 1 // m_totalLaps
	offset += 2 // m_trackLength
	
	sessionType := data[offset]
	info.SessionType = sessionTypes[sessionType]
	if info.SessionType == "" {
		info.SessionType = "Unknown"
	}
	offset += 1
	
	trackID := int8(data[offset])
	info.TrackName = trackNames[trackID]
	if info.TrackName == "" {
		info.TrackName = fmt.Sprintf("Track%d", trackID)
	}
	
	// Extract time of day (further in the packet)
	// Skip to timeOfDay offset (check documentation for exact position)
	// For now, we'll skip this as it's deep in the structure
}

// parseParticipantsPacket extracts player name from participants packet
func parseParticipantsPacket(data []byte, info *SessionInfo) {
	if len(data) < 30 {
		return
	}
	
	// Get player car index from header
	playerCarIndex := data[27]
	
	// Skip header (29 bytes)
	offset := 29
	
	// Read numActiveCars
	if offset >= len(data) {
		return
	}
	offset += 1 // m_numActiveCars
	
	// Each ParticipantData is approximately 58 bytes
	// Calculate offset to player's data
	participantOffset := offset + (int(playerCarIndex) * 58)
	
	if participantOffset+58 > len(data) {
		return
	}
	
	// Skip to name field (48 bytes into ParticipantData)
	nameOffset := participantOffset + 48
	
	if nameOffset+32 > len(data) {
		return
	}
	
	// Extract name (32 bytes, null-terminated UTF-8)
	nameBytes := data[nameOffset : nameOffset+32]
	name := extractNullTerminatedString(nameBytes)
	
	// Clean up the name
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Player"
	}
	
	// Remove any non-ASCII characters for filename safety
	name = sanitizeForFilename(name)
	
	info.PlayerName = name
}

// extractNullTerminatedString extracts a null-terminated string
func extractNullTerminatedString(data []byte) string {
	for i, b := range data {
		if b == 0 {
			return string(data[:i])
		}
	}
	return string(data)
}

// sanitizeForFilename removes characters that aren't safe for filenames
func sanitizeForFilename(s string) string {
	// Replace spaces with underscores
	s = strings.ReplaceAll(s, " ", "_")
	
	// Keep only alphanumeric, underscore, and dash
	result := strings.Builder{}
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// GenerateFilename creates a descriptive filename from session info
func (si *SessionInfo) GenerateFilename() string {
	parts := []string{}
	
	if si.TrackName != "" {
		parts = append(parts, si.TrackName)
	}
	
	if si.SessionType != "" && si.SessionType != "Unknown" {
		parts = append(parts, si.SessionType)
	}
	
	if si.PlayerName != "" && si.PlayerName != "Player" {
		parts = append(parts, si.PlayerName)
	}
	
	if si.Weather != "" && si.Weather != "Clear" {
		parts = append(parts, si.Weather)
	}
	
	if len(parts) == 0 {
		return "session"
	}
	
	return strings.Join(parts, "_")
}

// String returns a human-readable description
func (si *SessionInfo) String() string {
	if !si.HasInfo {
		return "No session info available"
	}
	
	parts := []string{}
	
	if si.PlayerName != "" {
		parts = append(parts, fmt.Sprintf("Player: %s", si.PlayerName))
	}
	if si.TrackName != "" {
		parts = append(parts, fmt.Sprintf("Track: %s", si.TrackName))
	}
	if si.SessionType != "" {
		parts = append(parts, fmt.Sprintf("Session: %s", si.SessionType))
	}
	if si.Weather != "" {
		parts = append(parts, fmt.Sprintf("Weather: %s", si.Weather))
	}
	
	return strings.Join(parts, " | ")
}
