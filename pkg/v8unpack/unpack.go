package v8unpack

import (
	"encoding/binary"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf16"
)

type v8address uint32

type Reader interface {
	ReadFragment(v8address, v8address) []byte
}

type FileReader struct {
	file *os.File
}

type BytesReader struct {
	data []byte
}

func (reader *FileReader) ReadFragment(begin v8address, length v8address) []byte {
	reader.file.Seek(int64(begin), 0)
	bufLength := v8address(length)

	buf := make([]byte, bufLength)
	for i := v8address(1); true; i++ {
		n, err := reader.file.Read(buf)

		if n > 0 {
			return buf
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			println("ReadFileFragment. ", err.Error())
			break
		}
	}

	return buf
}

func (reader *BytesReader) ReadFragment(begin v8address, length v8address) []byte {
	buf := reader.data[begin : begin+length]
	return buf
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

func convertFilename(filenameUTF16 []byte) string {
	utf := make([]uint16, len(filenameUTF16)/2)
	for i := 0; i < len(filenameUTF16); i += 2 {
		utf[(i / 2)] = binary.LittleEndian.Uint16(filenameUTF16[i:])
	}

	filename := string(utf16.Decode(utf))
	filename = strings.TrimRight(filename, string([]byte{0, 0}))
	return filename
}

func convertAddress(b []byte) v8address {
	bytes := ""
	for _, v := range b {
		bytes += string(v)
	}
	i, err := strconv.ParseUint(bytes, 16, 32)
	if err != nil {
		return 0
	}
	return v8address(i)
}
