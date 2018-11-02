//go:generate stringer -type=PacketType

package multiplexer

import (
	"encoding/hex"
	"errors"
)

// PacketType defines the packet type.
type PacketType byte

// Available packet types
const (
	PushData PacketType = iota
	PushACK
	PullData
	PullResp
	PullACK
	TXACK
)

// Protocol versions
const (
	ProtocolVersion1 uint8 = 0x01
	ProtocolVersion2 uint8 = 0x02
)

// GetPacketType returns the packet type for the given packet data.
func GetPacketType(data []byte) (PacketType, error) {
	if len(data) < 4 {
		return PacketType(0), errors.New("at least 4 bytes of data are expected")
	}
	return PacketType(data[3]), nil
}

// GetGatewayID returns the gateway ID for the given packet data.
func GetGatewayID(data []byte) (string, error) {
	if len(data) < 12 {
		return "", errors.New("at least 12 bytes of data are expected")
	}
	return hex.EncodeToString(data[4:12]), nil
}
