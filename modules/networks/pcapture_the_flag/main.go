package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"pcap/pkg/ethernet"
	"pcap/pkg/ipv4"
	"pcap/pkg/parse"
	"pcap/pkg/pcap"
	"pcap/pkg/tcp"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	f, err := os.Open("net.cap")
	must(err)
	dump, err := ioutil.ReadAll(f)
	must(err)
	parser := &parse.Parser{Raw: dump}

	_ = pcap.ReadFileHeader(parser)
	// fmt.Printf("PCAP file header: %#v\n", fileHeader)

	var totalPackets int
	msgs := []msg{}
	for parser.Idx < len(parser.Raw) {
		packet := pcap.ReadPacket(parser)
		totalPackets++

		// fmt.Printf("PCAP Packet headers: %s\n", packet.String())
		eParser := &parse.Parser{Raw: packet.Payload, OtherEndian: true}
		ethernet.ReadPacketHeaders(eParser)
		// fmt.Printf("Ethernet packet headers: %s\n", ePacket.String())

		ipStart := eParser.Idx
		ipPacket := ipv4.ReadPacketHeaders(eParser)

		if !ipPacket.V4() {
			panic("not ipv4")
		}
		if !ipPacket.TCP() {
			panic("not tcp")
		}

		tcpStart := eParser.Idx
		tcpPacket := tcp.ReadPacketHeaders(eParser)

		// 1 row == 32-bits. DataOffset = 6
		//    _____________
		//  0 |________} <- IP start
		//  1 |________}
		//  2 |________} <- TCP start
		//  3 |________}
		//  4 |________} <-- ME (Idx = 4)
		//  5 |________} <-- TCP OPTIONS
		//  6 |________}
		//  7 |________} <-- TCP DATA START (Idx = 6)
		//  8 |________}
		//  9 |________} <-- END (Idx = 9)

		target := tcpStart + tcpPacket.HeaderLengthBytes() - 1
		eParser.Read(target - eParser.Idx)
		if eParser.Idx != tcpStart+tcpPacket.HeaderLengthBytes()-1 {
			fmt.Printf("TCP Start: %d\nHeader Length: %d\nPos: %d\n", tcpStart, tcpPacket.HeaderLengthBytes(), eParser.Idx)
			panic("Not at TCP data")
		}
		target = ipStart + int(ipPacket.TotalLength) - 1
		eParser.Read(target - eParser.Idx)
		if eParser.Idx != ipStart+int(ipPacket.TotalLength)-1 {
			fmt.Printf("IP Start: %d\nTotal Length: %d\nPos: %d\n", ipStart, ipPacket.TotalLength, eParser.Idx)
			panic("Not at end of IP datagram")
		}
		if tcpPacket.SourcePort == 80 {
			// TODO(cjea): de-dup ?
			msgs = append(msgs, msg{
				Pos:     tcpPacket.SequenceNumber,
				Payload: eParser.Workspace,
			})
		}
	}
	s := []byte{}
	msgs = sortMsgs(msgs)
	for _, el := range msgs {
		s = append(s, el.Payload...)
	}
	var httpBodyIdx int
	for i := 0; i+3 < len(s); i++ {
		if s[i] == 0x0d && s[i+1] == 0x0a && s[i+2] == 0x0d && s[i+3] == 0x0a {
			httpBodyIdx = i + 4
			break
		}
	}
	fmt.Printf("%s", string(s[httpBodyIdx:]))
}

type msg struct {
	Pos     uint32
	Payload []byte
}

func hasSuffix(as []byte, bs []byte) bool {
	if len(bs) > len(as) {
		return false
	}
	slice := as[len(as)-len(bs):]
	for i := range slice {
		if slice[i] != bs[i] {
			return false
		}
	}
	return true
}

func sortMsgs(msgs []msg) []msg {
	length := len(msgs)
	if length < 2 {
		return msgs
	}
	mid := length / 2
	l := sortMsgs(msgs[0:mid])
	r := sortMsgs(msgs[mid:])
	count := 0
	ret := make([]msg, len(msgs))
	for ; count < length; count++ {
		if len(l) > 0 {
			if len(r) == 0 || l[0].Pos <= r[0].Pos {
				ret[count] = l[0]
				l = l[1:]
				continue
			} else {
				ret[count] = r[0]
				r = r[1:]
				continue
			}
		} else {
			if len(r) > 0 {
				ret[count] = r[0]
				r = r[1:]
				continue
			} else {
				return ret
			}
		}
	}
	return ret
}
