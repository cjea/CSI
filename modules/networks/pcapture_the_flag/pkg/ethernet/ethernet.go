package ethernet

import (
	"fmt"
	"pcap/pkg/parse"
)

type EthernetPacket struct {
	Destination []byte
	Source      []byte
	Type        uint16
}

func (p EthernetPacket) String() string {
	return fmt.Sprintf(
		"{\n\tDestination: %x,\n\tSource: %x,\n\tType: %x,\n\t",
		p.Destination, p.Source, p.Type,
	)
}

func ReadPacketHeaders(p *parse.Parser) EthernetPacket {
	ret := EthernetPacket{}
	p.Read(6)
	ret.Destination = p.Load()

	p.Read(6)
	ret.Source = p.Load()

	// Maybe the optional 802.1Q tag? If things are messed up, try removing this.
	// p.Read(4)

	ret.Type = p.ReadUInt16()
	return ret
}
