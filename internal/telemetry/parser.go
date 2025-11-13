package telemetry

import (
	"encoding/binary"
	"math"
)

// TelemetryData holds parsed telemetry values for display
type TelemetryData struct {
	Speed            float32
	Throttle         float32
	Brake            float32
	Gear             int8
	EngineRPM        uint16
	DRS              uint8
	EngineTemp       uint16
	TyreTemp         [4]uint8 // FL, FR, RL, RR
	TyrePressure     [4]float32
	FuelLevel        float32
	ERSStoreEnergy   float32
	ERSDeployMode    uint8
}

// ParseCarTelemetryPacket extracts telemetry data from packet ID 6
func ParseCarTelemetryPacket(data []byte, playerCarIndex uint8) *TelemetryData {
	if len(data) < 29 || data[6] != 6 { // Check if it's telemetry packet
		return nil
	}

	// Each car telemetry data is 60 bytes
	// Offset: 29 (header) + (playerCarIndex * 60)
	offset := 29 + (int(playerCarIndex) * 60)
	
	if offset+60 > len(data) {
		return nil
	}

	td := &TelemetryData{}
	
	// Parse telemetry data according to F1 25 spec
	// uint16 m_speed (km/h)
	td.Speed = float32(binary.LittleEndian.Uint16(data[offset : offset+2]))
	offset += 2
	
	// float m_throttle (0.0 to 1.0)
	td.Throttle = math.Float32frombits(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	
	// float m_steer (skip)
	offset += 4
	
	// float m_brake (0.0 to 1.0)
	td.Brake = math.Float32frombits(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	
	// uint8 m_clutch (skip)
	offset += 1
	
	// int8 m_gear
	td.Gear = int8(data[offset])
	offset += 1
	
	// uint16 m_engineRPM
	td.EngineRPM = binary.LittleEndian.Uint16(data[offset : offset+2])
	offset += 2
	
	// uint8 m_drs (0=off, 1=on)
	td.DRS = data[offset]
	offset += 1
	
	// uint8 m_revLightsPercent (skip)
	offset += 1
	
	// uint16 m_revLightsBitValue (skip)
	offset += 2
	
	// uint16 m_brakesTemperature[4] (skip)
	offset += 8
	
	// uint8 m_tyresSurfaceTemperature[4]
	for i := 0; i < 4; i++ {
		td.TyreTemp[i] = data[offset]
		offset += 1
	}
	
	// uint8 m_tyresInnerTemperature[4] (skip)
	offset += 4
	
	// uint16 m_engineTemperature
	td.EngineTemp = binary.LittleEndian.Uint16(data[offset : offset+2])
	offset += 2
	
	// float m_tyresPressure[4]
	for i := 0; i < 4; i++ {
		td.TyrePressure[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[offset : offset+4]))
		offset += 4
	}
	
	// uint8 m_surfaceType[4] (skip)
	
	return td
}

// ParseCarStatusPacket extracts status data from packet ID 7
func ParseCarStatusPacket(data []byte, playerCarIndex uint8) *TelemetryData {
	if len(data) < 29 || data[6] != 7 {
		return nil
	}

	// Each car status data is 58 bytes according to F1 25 spec
	offset := 29 + (int(playerCarIndex) * 58)
	
	if offset+58 > len(data) {
		return nil
	}

	td := &TelemetryData{}
	
	// uint8 m_tractionControl
	offset += 1
	
	// uint8 m_antiLockBrakes
	offset += 1
	
	// uint8 m_fuelMix
	offset += 1
	
	// uint8 m_frontBrakeBias
	offset += 1
	
	// uint8 m_pitLimiterStatus
	offset += 1
	
	// float m_fuelInTank
	td.FuelLevel = math.Float32frombits(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	
	// float m_fuelCapacity (skip)
	offset += 4
	
	// float m_fuelRemainingLaps (skip)
	offset += 4
	
	// uint16 m_maxRPM (skip)
	offset += 2
	
	// uint16 m_idleRPM (skip)
	offset += 2
	
	// uint8 m_maxGears (skip)
	offset += 1
	
	// uint8 m_drsAllowed (0 = not allowed, 1 = allowed)
	td.DRS = data[offset]
	offset += 1
	
	// uint16 m_drsActivationDistance (skip)
	offset += 2
	
	// uint8 m_actualTyreCompound (skip)
	offset += 1
	
	// uint8 m_visualTyreCompound (skip)
	offset += 1
	
	// uint8 m_tyresAgeLaps (skip)
	offset += 1
	
	// int8 m_vehicleFiaFlags (skip)
	offset += 1
	
	// float m_enginePowerICE (skip)
	offset += 4
	
	// float m_enginePowerMGUK (skip)
	offset += 4
	
	// float m_ersStoreEnergy
	td.ERSStoreEnergy = math.Float32frombits(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	
	// uint8 m_ersDeployMode
	td.ERSDeployMode = data[offset]
	
	return td
}

// MergeTelemetryData merges telemetry from multiple packet types
func MergeTelemetryData(base, new *TelemetryData) *TelemetryData {
	if base == nil {
		if new == nil {
			return &TelemetryData{}
		}
		return new
	}
	if new == nil {
		return base
	}
	
	// Update telemetry values from Car Telemetry packet (ID 6)
	// Only update if new value seems valid (Speed, RPM, Temp should be non-zero during driving)
	if new.Speed > 0 {
		base.Speed = new.Speed
		// Throttle and Brake come from same packet, update together
		base.Throttle = new.Throttle
		base.Brake = new.Brake
		base.Gear = new.Gear
		base.EngineRPM = new.EngineRPM
	}
	
	if new.EngineTemp > 0 {
		base.EngineTemp = new.EngineTemp
	}
	
	// Update tyre temps if any are non-zero
	if new.TyreTemp[0] > 0 || new.TyreTemp[1] > 0 || new.TyreTemp[2] > 0 || new.TyreTemp[3] > 0 {
		base.TyreTemp = new.TyreTemp
	}
	
	// Update values from Car Status packet (ID 7)
	if new.FuelLevel > 0 {
		base.FuelLevel = new.FuelLevel
		// ERS comes from same packet
		base.ERSStoreEnergy = new.ERSStoreEnergy
		base.DRS = new.DRS
	}
	
	return base
}
