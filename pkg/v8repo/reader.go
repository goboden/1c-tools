package v8repo

import (
	"io"
	"os"
)

type Reader interface {
	ReadFragment(uint, uint) []byte
}

type FileReader struct {
	file *os.File
}

func (r *FileReader) ReadFragment(begin uint, length uint) []byte {
	r.file.Seek(int64(begin), 0)
	bufLength := length

	buf := make([]byte, bufLength)
	for i := 1; true; i++ {
		n, err := r.file.Read(buf)

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

func NewFileReader(file *os.File) *FileReader {
	reader := &FileReader{file: file}
	return reader
}
