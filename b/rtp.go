package b

const (
	VERSION   = 2
	PADDING   = 0
	Extension = 0
	CC        = 0
	Marker    = 0
	Ssrc      = 1337

	HEADER_SIZE = 12
)

type RTPPacket struct {
	Header  []byte
	Payload []byte
}

func NewRTPPacket(ptType int, framenb int, time int, data []byte, dataLength int) *RTPPacket {
	packet := new(RTPPacket)
	packet.Header = make([]byte, HEADER_SIZE)
	packet.Header[0] = VERSION<<6 | PADDING<<5 | Extension<<4 | CC
	packet.Header[1] = Marker<<7 | byte(ptType&0x000000FF)
	packet.Header[2] = byte(framenb >> 8)
	packet.Header[3] = byte(framenb & 0xFF)
	packet.Header[4] = byte(time >> 24)
	packet.Header[5] = byte(time >> 16)
	packet.Header[6] = byte(time >> 8)
	packet.Header[7] = byte(time & 0xFF)
	packet.Header[8] = (byte)(Ssrc >> 24)
	packet.Header[9] = (byte)(Ssrc >> 16)
	packet.Header[10] = (byte)(Ssrc >> 8)
	packet.Header[11] = (byte)(Ssrc & 0xFF)

	packet.Payload = make([]byte, dataLength)
	copy(packet.Payload, data)

	return packet
}
