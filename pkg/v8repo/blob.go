package v8repo

import (
	"bytes"
	"encoding/binary"
)

const BlockLength = 256

type Blob struct {
	Lang   string
	Blocks []uint32
	Raw    []byte
	Data   [][]byte
}

func NewBlob(data []byte) *Blob {
	blob := &Blob{Raw: data}

	header := blob.ReadBlock(1)
	blob.Lang = string(header[:32])

	var numblocks uint32
	binary.Read(bytes.NewReader(header[32:36]), binary.LittleEndian, &numblocks)

	for i := 0; i < int(numblocks); i++ {
		var numblock uint32
		next := 36 + (i * 4)
		binary.Read(bytes.NewReader(header[next:next+4]), binary.LittleEndian, &numblock)
		blob.Blocks = append(blob.Blocks, numblock)
	}

	for _, n := range blob.Blocks {
		block := blob.ReadBlock(uint(n))
		blob.Data = append(blob.Data, block)
		// println("_", n)
		// fmt.Printf("%s\n", block)
	}

	return blob
}

func (b *Blob) ReadBlock(number uint) []byte {
	begin := BlockLength * number
	data := b.Raw[begin : begin+BlockLength]

	next := uint32(0)
	length := uint16(0)
	binary.Read(bytes.NewReader(data[0:4]), binary.LittleEndian, &next)
	binary.Read(bytes.NewReader(data[4:6]), binary.LittleEndian, &length)

	block := data[6 : 6+length]

	if next != 0 {
		nd := b.ReadBlock(uint(next))
		block = append(block, nd...)
	}

	return block
}
