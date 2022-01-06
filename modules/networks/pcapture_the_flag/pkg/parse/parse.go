package parse

import (
	"fmt"
	"unsafe"
)

// Read raw bytes, but flip the byte order of every processed chunk.
type Parser struct {
	Raw         []byte
	Workspace   []byte
	Idx         int
	OtherEndian bool
}

func (p *Parser) Read(numBytes int) {
	if p.Idx+numBytes > len(p.Raw) {
		fmt.Printf(
			"can't read %d bytes (already at %d of %d)\n",
			numBytes, p.Idx, len(p.Raw),
		)
		panic("Read past the end of the pcap dump")
	}

	idx := p.Idx
	p.Workspace = p.Raw[idx : idx+numBytes]
	p.Idx += numBytes
}

func (p *Parser) Load() []byte {
	last := len(p.Workspace) - 1
	ret := make([]byte, last+1)
	for i := range p.Workspace {
		idx := i
		if p.OtherEndian {
			idx = last - i
		}
		ret[i] = p.Workspace[idx]
	}
	return ret
}

func (p *Parser) ReadOneByte() byte {
	p.Read(1)
	return p.Load()[0]
}

func (p *Parser) ReadInt32() int32 {
	p.Read(4)
	return *unsafeI32(p.Load())
}

func (p *Parser) ReadInt16() int16 {
	p.Read(2)
	return *unsafeI16(p.Load())
}

func (p *Parser) ReadUInt32() uint32 {
	p.Read(4)
	return *unsafeU32(p.Load())
}

func (p *Parser) ReadUInt16() uint16 {
	p.Read(2)
	return *unsafeU16(p.Load())
}

func unsafeI16(bs []byte) *int16 {
	if len(bs) != 2 {
		panic("int16 requires 2 bytes")
	}

	ptr := unsafe.Pointer(&bs[0])
	return (*int16)(ptr)
}

func unsafeI32(bs []byte) *int32 {
	if len(bs) != 4 {
		panic("int32 requires 4 bytes")
	}

	ptr := unsafe.Pointer(&bs[0])
	return (*int32)(ptr)
}

func unsafeU16(bs []byte) *uint16 {
	if len(bs) != 2 {
		panic("uint16 requires 2 bytes")
	}

	ptr := unsafe.Pointer(&bs[0])
	return (*uint16)(ptr)
}

func unsafeU32(bs []byte) *uint32 {
	if len(bs) != 4 {
		panic("uint32 requires 4 bytes")
	}

	ptr := unsafe.Pointer(&bs[0])
	return (*uint32)(ptr)
}

func reverse(bs []byte) {
	last := len(bs) - 1
	for i := range bs {
		tmp := bs[i]
		bs[i] = bs[last-i]
		bs[last-i] = tmp
	}
}
