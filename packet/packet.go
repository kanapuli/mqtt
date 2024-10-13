package packet

import (
	"errors"
	"fmt"
)

const (
	MAX_REMAINING_LENGTH = 128 * 128 * 128 * 128
)

var (
	ErrInsufficientFixedHeader            = errors.New("insufficient data for fixed header")
	ErrMalformedRemainingLength           = errors.New("malformed Remaining Length in the fixed header")
	ErrInsufficientDataForRemainingLength = errors.New("insufficient data for Remaining Length")
)

// ControlPacket represents the MQTT control packet.
type ControlPacket struct {
	FixedHeader    any
	VariableHeader any
	Payload        any
}

// FixedHeader has the information about the control type and the remaining length of bytes
// including the variable headers and the payload
//
// Representation of FixedHeader
// Bit		|	7	|	6	|	5	|	4	|	3	|	2	|	1	|	0	|
// byte1	|  Mqtt control packet type		| Flags to each ctrl pkt type	|
// byte2 ..	|  Remaining Length of bytes inside the ctrl pkt				|
type FixedHeader struct {
	PacketType byte
	Flags      byte
	// RemainingLength represents the number of bytes remaining within the current Control packet,
	// including the data in the Variable header and the Payload. The remaining length does not
	// include the bytes used to encode the remaining length
	// The data representation of the RemainingLength is variable byte integer
	RemainingLength int
}

// NewFixedHeader creates a new fixed header for the control packets
func NewFixedHeader(packetTypeAndFlags byte, remainingLen int) (*FixedHeader, error) {
	packetType := shiftNBits(4, packetTypeAndFlags)
	flags := getLowerFourBits(packetTypeAndFlags)
	return &FixedHeader{
		PacketType:      packetType,
		Flags:           flags,
		RemainingLength: remainingLen,
	}, nil
}

func (fh *FixedHeader) String() string {
	return fmt.Sprintf("PacketType: %d, Flags: %04b, RemainingLength: %d", fh.PacketType, fh.Flags, fh.RemainingLength)
}

func parseFixedHeader(data []byte) (*FixedHeader, error) {
	if len(data) < 2 {
		return nil, ErrInsufficientFixedHeader
	}
	packetTypeAndFlags := data[0]
	remainingLen := 0
	remainingLen, _, err := decodeVariableByteInteger(data[1:])
	if err != nil {
		return nil, err
	}
	return NewFixedHeader(packetTypeAndFlags, remainingLen)
}

func decodeVariableByteInteger(data []byte) (int, int, error) {
	var (
		value      int
		multiplier int = 1
		bytesRead  int
	)
	const maxMultiplierValue = 128 * 128 * 128

	for _, b := range data {
		bytesRead++
		// Mask out the MSB which is the continuation bit
		digit := int(b & 0x7F)
		value += digit * multiplier

		// Return the value if the continuation bit is set to 0
		if b&0x80 == 0 {
			return value, bytesRead, nil
		}

		multiplier *= 128
		if multiplier > maxMultiplierValue {
			return 0, 0, ErrMalformedRemainingLength
		}

	}
	return 0, 0, ErrInsufficientDataForRemainingLength
}

// shiftNBits shifts a byte 'd' right by 'n' bits
func shiftNBits(n int, d byte) byte {
	return d >> n
}

// getLowerFourBits return the least significant 4 bits
func getLowerFourBits(d byte) byte {
	return d & 0x0F
}
