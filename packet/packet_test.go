package packet

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
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
		expectedValue int
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
		input    int
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
		{
			name:     "Negative Value",
			input:    -1,
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
