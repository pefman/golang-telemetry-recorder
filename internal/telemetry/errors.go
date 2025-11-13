package telemetry

import "errors"

var (
	ErrInvalidPacket = errors.New("invalid packet data")
	ErrTimeout       = errors.New("packet receive timeout")
	ErrStopped       = errors.New("receiver stopped")
)
