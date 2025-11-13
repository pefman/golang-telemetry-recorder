package telemetry

import (
	"encoding/binary"
	"time"
	"unsafe"
)

// PacketType represents different F1 telemetry packet types
type PacketType uint8

const (
	PacketMotion PacketType = iota
	PacketSession
	PacketLapData
	PacketEvent
	PacketParticipants
	PacketCarSetups
	PacketCarTelemetry
	PacketCarStatus
	PacketFinalClassification
	PacketLobbyInfo
	PacketCarDamage
	PacketSessionHistory
	PacketTyreSets
	PacketMotionEx
)

// PacketHeader is the common header for all F1 telemetry packets
type PacketHeader struct {
	PacketFormat            uint16
	GameYear                uint8
	GameMajorVersion        uint8
	GameMinorVersion        uint8
	PacketVersion           uint8
	PacketID                uint8
	SessionUID              uint64
	SessionTime             float32
	FrameIdentifier         uint32
	OverallFrameIdentifier  uint32
	PlayerCarIndex          uint8
	SecondaryPlayerCarIndex uint8
}

// RecordedPacket represents a packet with timestamp for recording/playback
type RecordedPacket struct {
	Timestamp time.Time
	Data      []byte
	Header    PacketHeader
}

// ParseHeader extracts the packet header from raw data
func ParseHeader(data []byte) (*PacketHeader, error) {
	if len(data) < 29 { // Minimum header size
		return nil, ErrInvalidPacket
	}

	header := &PacketHeader{
		PacketFormat:            binary.LittleEndian.Uint16(data[0:2]),
		GameYear:                data[2],
		GameMajorVersion:        data[3],
		GameMinorVersion:        data[4],
		PacketVersion:           data[5],
		PacketID:                data[6],
		SessionUID:              binary.LittleEndian.Uint64(data[7:15]),
		SessionTime:             float32FromBytes(data[15:19]),
		FrameIdentifier:         binary.LittleEndian.Uint32(data[19:23]),
		OverallFrameIdentifier:  binary.LittleEndian.Uint32(data[23:27]),
		PlayerCarIndex:          data[27],
		SecondaryPlayerCarIndex: data[28],
	}

	return header, nil
}

// float32FromBytes converts 4 bytes to float32
func float32FromBytes(b []byte) float32 {
	bits := binary.LittleEndian.Uint32(b)
	return float32frombits(bits)
}

// float32frombits converts uint32 to float32
func float32frombits(b uint32) float32 {
	return *(*float32)(unsafe.Pointer(&b))
}

// GetPacketTypeName returns human-readable packet type name
func GetPacketTypeName(packetID uint8) string {
	names := map[uint8]string{
		0:  "Motion",
		1:  "Session",
		2:  "Lap Data",
		3:  "Event",
		4:  "Participants",
		5:  "Car Setups",
		6:  "Car Telemetry",
		7:  "Car Status",
		8:  "Final Classification",
		9:  "Lobby Info",
		10: "Car Damage",
		11: "Session History",
		12: "Tyre Sets",
		13: "Motion Ex",
	}

	if name, ok := names[packetID]; ok {
		return name
	}
	return "Unknown"
}
