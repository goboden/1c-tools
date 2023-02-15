package v8unpack

import (
	"encoding/binary"
	"os"
	"strconv"
	"strings"
	"unicode/utf16"
)

type v8address uint32

type Reader interface {
	ReadFragment(v8address, v8address) ([]byte, bool)
}

type FileReader struct {
	file *os.File
}

type BytesReader struct {
	data []byte
}

func (reader *FileReader) ReadFragment(begin v8address, length v8address) ([]byte, bool) {
	buf := make([]byte, v8address(length))
	reader.file.Seek(int64(begin), 0)

	n, err := reader.file.Read(buf)
	if err != nil {
		return nil, false
	}

	return buf[:n], true
}

func (reader *BytesReader) ReadFragment(begin v8address, length v8address) ([]byte, bool) {
	datalen := v8address(len(reader.data))
	if begin > datalen {
		return nil, false
	}

	end := begin + length
	if end > datalen {
		end = datalen
	}

	buf := reader.data[begin:end]
	return buf, true
}

func NewFileReader(file *os.File) *FileReader {
	reader := new(FileReader)
	reader.file = file
	return reader
}

func NewBytesReader(data []byte) *BytesReader {
	reader := new(BytesReader)
	reader.data = data
	return reader
}

func bytesToString(b []byte) string {
	utf := make([]uint16, len(b)/2)
	for i := 0; i < len(b); i += 2 {
		utf[(i / 2)] = binary.LittleEndian.Uint16(b[i:])
	}

	str := string(utf16.Decode(utf))
	str = strings.TrimRight(str, string([]byte{0, 0}))
	return str
}

func bytesToAdr(b []byte) v8address {
	var adr string

	for _, v := range b {
		adr += string(v)
	}
	i, err := strconv.ParseUint(adr, 16, 32)
	if err != nil {
		return 0
	}
	return v8address(i)
}

// func parseAddr(b []byte)
