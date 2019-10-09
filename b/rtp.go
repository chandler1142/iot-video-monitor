package b

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	VERSION        = 2    //V: RFC3550中规定的版本号为2
	PADDING        = 0    //P: 如设置填充位，在包末尾包含了额外的附加信息
	Extension      = 0    //X: 如果该位被设置，则在固定的头部后存在一个扩展头部
	CC             = 0    //CC: CSRC计数包括紧接在固定头后标识CSRC个数
	Marker         = 0    //M: 一般是0, 对于H264，如果当前 NALU为一个接入单元最后的那个NALU，那么将M位置 1
	Ssrc           = 1337 // 使同一RTP连接中两个同步源传输的数据包中没有相同的SSRC标识
	HEADER_SIZE    = 12
	MTU            = 1500
	TS_PACKET_SIZE = 188
)

type RTPPacket struct {
	Header  []byte
	Payload []byte
}

/**
	ptType: RTP packet所携带信息的类型，标准类型列出在RFC3551中

 */
func NewRTPPacket(ptType int, sequenceNumber int, timestamp int, data []byte, dataLength int) *RTPPacket {
	packet := new(RTPPacket)
	packet.Header = make([]byte, HEADER_SIZE)
	packet.Header[0] = VERSION<<6 | PADDING<<5 | Extension<<4 | CC
	packet.Header[1] = Marker<<7 | byte(ptType&0x000000FF)
	packet.Header[2] = byte(sequenceNumber >> 8)
	packet.Header[3] = byte(sequenceNumber & 0xFF)
	packet.Header[4] = byte(timestamp >> 24)
	packet.Header[5] = byte(timestamp >> 16)
	packet.Header[6] = byte(timestamp >> 8)
	packet.Header[7] = byte(timestamp & 0xFF)
	packet.Header[8] = (byte)(Ssrc >> 24)
	packet.Header[9] = (byte)(Ssrc >> 16)
	packet.Header[10] = (byte)(Ssrc >> 8)
	packet.Header[11] = (byte)(Ssrc & 0xFF)

	packet.Payload = make([]byte, dataLength)
	copy(packet.Payload, data)

	return packet
}

func NewEmptyRTPPacket(ptType int, sequenceNumber int, timestamp int) *RTPPacket {
	packet := new(RTPPacket)
	packet.Header = make([]byte, HEADER_SIZE)
	packet.Header[0] = VERSION<<6 | PADDING<<5 | Extension<<4 | CC
	packet.Header[1] = Marker<<7 | byte(ptType&0x000000FF)
	packet.Header[2] = byte(sequenceNumber >> 8)
	packet.Header[3] = byte(sequenceNumber & 0xFF)
	packet.Header[4] = byte(timestamp >> 24)
	packet.Header[5] = byte(timestamp >> 16)
	packet.Header[6] = byte(timestamp >> 8)
	packet.Header[7] = byte(timestamp & 0xFF)
	packet.Header[8] = (byte)(Ssrc >> 24)
	packet.Header[9] = (byte)(Ssrc >> 16)
	packet.Header[10] = (byte)(Ssrc >> 8)
	packet.Header[11] = (byte)(Ssrc & 0xFF)

	return packet
}

func (packet *RTPPacket) SetData(data []byte, clear bool) {
	if clear {
		packet.Payload = []byte{}
	} else {
		packet.Payload = append(packet.Payload, data...)
	}
}

func (packet *RTPPacket) GetBytes() []byte {
	return append(packet.Header, packet.Payload...)
}

func (packet *RTPPacket) SequenceIncrement() {
	seqBytes := append([]byte{}, packet.Header[2], packet.Header[3])
	var seq uint16 = 0
	bytesBuffer := bytes.NewBuffer(seqBytes)
	err := binary.Read(bytesBuffer, binary.BigEndian, &seq)
	if err != nil {
		fmt.Println(err)
		return
	}
	newSeq := seq + 1
	fmt.Printf("Original seq: %d, newSeq: %d \n", seq, newSeq)
	packet.Header[2] = byte(newSeq >> 8)
	packet.Header[3] = byte(newSeq & 0xFF)
}
