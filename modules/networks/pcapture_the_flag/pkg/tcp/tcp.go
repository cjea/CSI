package tcp

import (
	"pcap/pkg/parse"
)

type Packet struct {
	SourcePort      uint16
	DestinationPort uint16
	SequenceNumber  uint32
	AckNumber       *uint32
	DataOffset      uint8
	Action          uint8
	WindowSize      uint16
	CheckSum        uint16
	UrgentPointer   uint16
}

func ReadPacketHeaders(p *parse.Parser) Packet {
	ret := Packet{}

	ret.SourcePort = p.ReadUInt16()
	ret.DestinationPort = p.ReadUInt16()
	ret.SequenceNumber = p.ReadUInt32()
	ackNum := p.ReadUInt32()
	ret.AckNumber = &ackNum

	ret.DataOffset = p.ReadOneByte() >> 4
	ret.Action = p.ReadOneByte()
	ret.WindowSize = p.ReadUInt16()
	ret.CheckSum = p.ReadUInt16()
	ret.UrgentPointer = p.ReadUInt16()

	return ret
}

func (p Packet) HeaderLengthBytes() int {
	return int(p.DataOffset) * 4
}
