package pcap

import (
	"fmt"
	"pcap/pkg/parse"
	"time"
)

type PcapPacket struct {
	Timestamp  time.Time
	NCaptured  uint32
	NPotential uint32
	Payload    []byte
}

func (packet PcapPacket) String() string {
	return fmt.Sprintf(
		"{\n\tTimestamp: %v,\n\tCaptured: %v,\n\tPotential: %v,\n\tPayload length: %v,\n}",
		packet.Timestamp, packet.NCaptured, packet.NPotential, len(packet.Payload),
	)
}

func ReadPacket(p *parse.Parser) PcapPacket {
	ret := PcapPacket{}

	ret.Timestamp = time.Unix(int64(p.ReadInt32()), 0)
	p.Read(4) // Microseconds of timestamp

	ret.NCaptured = p.ReadUInt32()
	ret.NPotential = p.ReadUInt32()
	p.Read(int(ret.NCaptured))
	ret.Payload = p.Load()

	return ret
}

func ReadFileHeader(p *parse.Parser) FileHeader {
	ret := FileHeader{}

	p.Read(4)
	ret.MagicNumber = p.Workspace // Keep the magic number as is, to signal endian-ness.
	ret.Major = p.ReadInt16()
	ret.Minor = p.ReadInt16()
	p.Read(8) // throw away time zone offset and time stamp accuracy
	ret.SnapshotLen = p.ReadUInt32()
	ret.LinkLayerType = p.ReadUInt32()

	return ret
}

type FileHeader struct {
	MagicNumber   []byte
	Major         int16
	Minor         int16
	SnapshotLen   uint32
	LinkLayerType uint32
}
