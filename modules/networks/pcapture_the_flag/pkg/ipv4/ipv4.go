package ipv4

import "pcap/pkg/parse"

type Datagram struct {
	Version        byte
	IHL            byte
	DCSP           byte
	TotalLength    uint16
	Identification uint16
	Flags          byte // bit 0: Reserved; must be zero. bit 1: Don't Fragment (DF) bit 2: More Fragments (MF)
	FragmentOffset uint16
	TTL            byte
	Protocol       byte
	HeaderChecksum uint16
	Source         uint32
	Destination    uint32
}

func ReadPacketHeaders(p *parse.Parser) Datagram {
	ret := Datagram{}
	var tmp byte

	tmp = p.ReadOneByte()
	ret.Version = tmp >> 4
	ret.IHL = tmp & 0xf

	ret.DCSP = p.ReadOneByte()

	ret.TotalLength = p.ReadUInt16()
	ret.Identification = p.ReadUInt16()

	i := p.ReadUInt16()
	ret.Flags = byte(i >> 13)
	ret.FragmentOffset = i & 0x1fff

	ret.TTL = p.ReadOneByte()
	ret.Protocol = p.ReadOneByte()

	ret.HeaderChecksum = p.ReadUInt16()
	ret.Source = p.ReadUInt32()
	ret.Destination = p.ReadUInt32()

	return ret
}

func (d Datagram) V4() bool {
	return d.Version == 4
}

// https://www.wikiwand.com/en/List_of_IP_protocol_numbers
func (d Datagram) TCP() bool {
	return d.Protocol == 0x6
}
