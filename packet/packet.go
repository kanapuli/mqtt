package packet

// ControlPacket represents the MQTT control packet.
type ControlPacket struct {
	FixedHeader    any
	VariableHeader any
	Payload        any
}

// FixedHeader is a 2 byte information present inside
// the control packet
//
// Representation of FixedHeader
// Bit	|	7	|	6	|	5	|	4	|	3	|	2	|	1	|	0	|
// byte1|  Mqtt control packet type		| Flags to each ctrl pkt type	|
// byte2|			Remaining Length 									|
type FixedHeader struct {
	PacketType      byte
	Flags           byte
	RemainingLength int
}

// NewFixedHeader creates a new fixed header for the control packets
func NewFixedHeader(packetTypeAndFlags byte, remainingLen int) (*FixedHeader, error) {
	packetType := shiftNBits(4, packetTypeAndFlags)
	flags := packetTypeAndFlags & 0x0F
	return &FixedHeader{
		PacketType:      packetType,
		Flags:           flags,
		RemainingLength: remainingLen,
	}, nil
}

// shiftNBits shifts a byte 'd' right by 'n' bits
func shiftNBits(n int, d byte) byte {
	return d >> n
}
