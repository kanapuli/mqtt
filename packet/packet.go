package packet

import (
	"errors"
	"fmt"
)

const MAX_VARIABLE_BYTE_INTEGER_SIZE = (128 * 128 * 128 * 128) - 1

var (
	ErrInsufficientFixedHeader            = errors.New("insufficient data for fixed header")
	ErrMalformedRemainingLength           = errors.New("malformed Remaining Length in the fixed header")
	ErrInsufficientDataForRemainingLength = errors.New("insufficient data for Remaining Length")
	ErrVarByteIntegerOutOfRange           = errors.New("value out of range for variable byte integer")
)

// Maximum lenght of a VariableByteInteger can be 4 bytes
// See: https://docs.oasis-open.org/mqtt/mqtt/v5.0/os/mqtt-v5.0-os.html#_Toc3901011
type VariableByteInteger = uint32

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
	RemainingLength VariableByteInteger
}

// NewFixedHeader creates a new fixed header for the control packets
func NewFixedHeader(packetTypeAndFlags byte, remainingLen VariableByteInteger) (*FixedHeader, error) {
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
	var remainingLen uint32
	remainingLen, _, err := decodeVariableByteInteger(data[1:])
	if err != nil {
		return nil, err
	}
	return NewFixedHeader(packetTypeAndFlags, remainingLen)
}

func encodeVariableByteInteger(x VariableByteInteger) ([]byte, error) {
	if x < 0 || x > MAX_VARIABLE_BYTE_INTEGER_SIZE {
		return nil, ErrVarByteIntegerOutOfRange
	}
	var encodedBytes []byte
	for {
		// Get the least significant 7 bits
		encodedByte := byte(x % 128)
		x /= 128
		if x > 0 {
			// set the continuation bit
			encodedByte |= 0x80
		}
		encodedBytes = append(encodedBytes, encodedByte)
		if x <= 0 {
			break
		}
	}
	return encodedBytes, nil
}

func decodeVariableByteInteger(data []byte) (VariableByteInteger, int, error) {
	var (
		value      VariableByteInteger
		multiplier uint32 = 1
		bytesRead  int
	)
	const maxMultiplierValue = 128 * 128 * 128

	for _, b := range data {
		bytesRead++
		// Mask out the MSB which is the continuation bit
		digit := uint32(b & 0x7F)
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
