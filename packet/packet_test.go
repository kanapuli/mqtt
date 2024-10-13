package packet

import (
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
