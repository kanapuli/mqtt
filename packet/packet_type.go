package packet

// ControlPacketType represents the type of the control packet
// and corresponds to the first 4 bits of the fixed header
type ControlPacketType byte

const (
	Reserved ControlPacketType = iota
	Connect
	ConnAck
	Publish
	PubAck
	PubRec
	PubRel
	PubComp
	Subscribe
	SubAck
	Unsubscribe
	UnsubAck
	PingReq
	PingResp
	Disconnect
	Auth
)

func (c ControlPacketType) String() string {
	switch c {
	case Connect:
		return "Connect"
	case ConnAck:
		return "ConnAck"
	case Publish:
		return "Publish"
	case PubAck:
		return "PubAck"
	case PubRec:
		return "PubRec"
	case PubRel:
		return "PubRel"
	case PubComp:
		return "PubComp"
	case Subscribe:
		return "Subscribe"
	case SubAck:
		return "SubAck"
	case Unsubscribe:
		return "Unsubscribe"
	case UnsubAck:
		return "UnsubAck"
	case PingReq:
		return "PingReq"
	case PingResp:
		return "PingResp"
	case Disconnect:
		return "Disconnect"
	case Auth:
		return "Auth"
	default:
		return "Reserved"
	}
}
