package packet

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShiftNBits(t *testing.T) {
	testCases := []struct {
		name           string
		packet         byte
		expectedPacket byte
		n              int
	}{
		{
			name:           "Shift 1 bits right",
			packet:         0b11110000,
			expectedPacket: 0b01111000,
			n:              1,
		},
		{
			name:           "Shift 2 bits right",
			packet:         0b11110000,
			expectedPacket: 0b00111100,
			n:              2,
		},
		{
			name:           "Shift 3 bits right",
			packet:         0b11110000,
			expectedPacket: 0b00011110,
			n:              3,
		},
		{
			name:           "Shift 4 bits right",
			packet:         0b11110000,
			expectedPacket: 0b1111,
			n:              4,
		},
		{
			name:           "Shift 5 bits right",
			packet:         0b11110000,
			expectedPacket: 0b0111,
			n:              5,
		},
		{
			name:           "Shift 6 bits right",
			packet:         0b11110000,
			expectedPacket: 0b0011,
			n:              6,
		},
		{
			name:           "Shift 7 bits right",
			packet:         0b11110000,
			expectedPacket: 0b0001,
			n:              7,
		},
		{
			name:           "Shift 8 bits right",
			packet:         0b11110000,
			expectedPacket: 0b0000,
			n:              8,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := shiftNBits(tc.n, tc.packet)
			assert.Equal(t, tc.expectedPacket, actual)
		})
	}
}

func TestGetLowerFourBits(t *testing.T) {
	testCases := []struct {
		input    byte
		expected byte
	}{
		{0b11110000, 0b0000}, // Upper 4 bits set, lower 4 bits are 0
		{0b10101010, 0b1010}, // Lower 4 bits are 1010
		{0b00001111, 0b1111}, // All lower 4 bits set
		{0b00000000, 0b0000}, // All bits are 0
		{0b11111111, 0b1111}, // All bits are set, lower 4 bits are 1111
	}

	for _, tc := range testCases {
		actual := getLowerFourBits(tc.input)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestDecodeRemainingLength(t *testing.T) {
	testCases := []struct {
		name          string
		input         []byte
		expectedValue uint32
		expectedBytes int
		expectedErr   error
	}{
		{
			name:          "Single byte - Max Value",
			input:         []byte{0x7F},
			expectedValue: 127,
			expectedBytes: 1,
			expectedErr:   nil,
		},
		{
			name:          "Two bytes - Example value 134",
			input:         []byte{0x86, 0x01},
			expectedValue: 134,
			expectedBytes: 2,
			expectedErr:   nil,
		},
		{
			name:          "Three bytes - Example value 16,383",
			input:         []byte{0xFF, 0x7F},
			expectedValue: 16383,
			expectedBytes: 2,
			expectedErr:   nil,
		},
		{
			name:          "Four bytes - Max Value",
			input:         []byte{0xFF, 0xFF, 0xFF, 0x7F},
			expectedValue: 268435455,
			expectedBytes: 4,
			expectedErr:   nil,
		},
		{
			name:          "Malformed - Value too large",
			input:         []byte{0xFF, 0xFF, 0xFF, 0x80},
			expectedValue: 0,
			expectedBytes: 0,
			expectedErr:   errors.New("malformed Remaining Length in the fixed header"),
		},
		{
			name:          "Insufficient data",
			input:         []byte{0x80, 0x80},
			expectedValue: 0,
			expectedBytes: 0,
			expectedErr:   errors.New("insufficient data for Remaining Length"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, bytesRead, err := decodeVariableByteInteger(tc.input)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedValue, value)
			assert.Equal(t, tc.expectedBytes, bytesRead)
		})
	}
}

func TestEncodeVariableByteInteger(t *testing.T) {
	testCases := []struct {
		name     string
		input    VariableByteInteger
		expected []byte
		err      error
	}{
		{
			name:     "Single byte - Max Value 127",
			input:    127,
			expected: []byte{0x7F},
			err:      nil,
		},
		{
			name:     "Two bytes - Value 128",
			input:    128,
			expected: []byte{0x80, 0x01},
			err:      nil,
		},
		{
			name:     "Two bytes - Value 16383",
			input:    16383,
			expected: []byte{0xFF, 0x7F},
			err:      nil,
		},
		{
			name:     "Three bytes - Value 2097151",
			input:    2097151,
			expected: []byte{0xFF, 0xFF, 0x7F},
			err:      nil,
		},
		{
			name:     "Four bytes - Max Value 268435455",
			input:    268435455,
			expected: []byte{0xFF, 0xFF, 0xFF, 0x7F},
			err:      nil,
		},
		{
			name:     "Out of range - Value 268435456",
			input:    268435456,
			expected: nil,
			err:      ErrVarByteIntegerOutOfRange,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := encodeVariableByteInteger(tc.input)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, encoded)
		})
	}
}

func TestControlPacketTypeValueForFixedHeader(t *testing.T) {
	testCases := []struct {
		name            string
		input           byte
		remainingLength VariableByteInteger
		expected        ControlPacketType
		err             error
	}{
		{
			name:            "Reserved",
			input:           0x00,
			remainingLength: 0,
			expected:        Reserved,
			err:             nil,
		},
		{
			name:            "Connect",
			input:           0x10,
			remainingLength: 0,
			expected:        Connect,
			err:             nil,
		},
		{
			name:            "ConnAck",
			input:           0x20,
			remainingLength: 0,
			expected:        ConnAck,
			err:             nil,
		},
		{
			name:            "Publish",
			input:           0x30,
			remainingLength: 0,
			expected:        Publish,
			err:             nil,
		},
		{
			name:            "PubAck",
			input:           0x40,
			remainingLength: 0,
			expected:        PubAck,
			err:             nil,
		},
		{
			name:            "PubRec",
			input:           0x50,
			remainingLength: 0,
			expected:        PubRec,
			err:             nil,
		},
		{
			name:            "PubRel",
			input:           0x60,
			remainingLength: 0,
			expected:        PubRel,
			err:             nil,
		},
		{
			name:            "PubComp",
			input:           0x70,
			remainingLength: 0,
			expected:        PubComp,
			err:             nil,
		},
		{
			name:            "Subscribe",
			input:           0x80,
			remainingLength: 0,
			expected:        Subscribe,
			err:             nil,
		},
		{
			name:            "SubAck",
			input:           0x90,
			remainingLength: 0,
			expected:        SubAck,
			err:             nil,
		},
		{
			name:            "Unsubscribe",
			input:           0xA0,
			remainingLength: 0,
			expected:        Unsubscribe,
			err:             nil,
		},
		{
			name:            "UnsubAck",
			input:           0xB0,
			remainingLength: 0,
			expected:        UnsubAck,
			err:             nil,
		},
		{
			name:            "PingReq",
			input:           0xC0,
			remainingLength: 0,
			expected:        PingReq,
			err:             nil,
		},
		{
			name:            "PingResp",
			input:           0xD0,
			remainingLength: 0,
			expected:        PingResp,
			err:             nil,
		},
		{
			name:            "Disconnect",
			input:           0xE0,
			remainingLength: 0,
			expected:        Disconnect,
			err:             nil,
		},
		{
			name:            "Auth",
			input:           0xF0,
			remainingLength: 0,
			expected:        Auth,
			err:             nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fh, err := NewFixedHeader(tc.input, tc.remainingLength)
			assert.NoError(t, err)
			actual := fh.ControlPacketType()
			assert.Equal(t, tc.expected, actual)
		})
	}

}
