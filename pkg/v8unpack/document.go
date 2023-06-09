package v8unpack

import (
	"bytes"
	"compress/flate"
	"fmt"
)

const (
	blockHeaderLength = 31
	intMax            = 2147483647
)

type v8blockHeader struct {
	DocumentLength v8address
	BlockLength    v8address
	NextBlock      v8address
}

func readDocument(reader Reader, begin v8address) []byte {
	documentData := make([]byte, 0)
	documentLength := v8address(0)

	next := begin
	for {
		header, data := readBlock(reader, next)

		if header.DocumentLength != 0 {
			documentLength = header.DocumentLength
		}

		if header.NextBlock == intMax {
			documentData = append(documentData, data[:(documentLength-v8address(len(documentData)))]...)
			break
		}
		documentData = append(documentData, data...)

		next = header.NextBlock
	}
	return documentData
}

func readBlock(reader Reader, begin v8address) (*v8blockHeader, []byte) {
	header, _ := reader.ReadFragment(begin, blockHeaderLength)

	if header[0] != 13 {
		panic(fmt.Sprintf("! %d", header[0]))
	}

	blockHeader := new(v8blockHeader)

	blockHeader.DocumentLength = bytesToAdr(header[2:10])
	blockHeader.BlockLength = bytesToAdr(header[11:19])
	blockHeader.NextBlock = bytesToAdr(header[20:28])

	data, _ := reader.ReadFragment(begin+blockHeaderLength, blockHeader.BlockLength)

	return blockHeader, data
}

func readContent(reader Reader, begin v8address, defl bool) []byte {
	data := readDocument(reader, begin)

	if defl {
		return deflate(data)
	}

	return data
}

func deflate(data []byte) []byte {
	reader := flate.NewReader(bytes.NewReader(data))
	buffer := make([]byte, 1024)
	out := make([]byte, 0)

	for {
		n, _ := reader.Read(buffer)
		if n < len(buffer) {
			out = append(out, buffer[:n]...)
			break
		}
		out = append(out, buffer...)
	}
	return out
}
