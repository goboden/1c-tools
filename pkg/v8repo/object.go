package v8repo

import (
	"bytes"
	"encoding/binary"
)

type Object struct {
	Database *Database
	Type     uint32 // 0xFD1C short || 0x01FD1C extended
	Version1 uint32
	Version2 uint32
	Version3 uint32
	Length   uint64
	Pages    []uint32
}

func NewObject(database *Database, page uint32) *Object {
	obj := &Object{Database: database}

	data, _ := obj.Database.ReadPage(uint(page))

	binary.Read(bytes.NewReader(data[0:4]), binary.LittleEndian, &obj.Type)
	binary.Read(bytes.NewReader(data[4:8]), binary.LittleEndian, &obj.Version1)
	binary.Read(bytes.NewReader(data[8:12]), binary.LittleEndian, &obj.Version2)
	binary.Read(bytes.NewReader(data[12:16]), binary.LittleEndian, &obj.Version3)
	binary.Read(bytes.NewReader(data[16:24]), binary.LittleEndian, &obj.Length)

	obj.Pages = make([]uint32, 0)
	var block uint32
	for i := 24; i <= len(data); i += 4 {
		binary.Read(bytes.NewReader(data[i:(i+4)]), binary.LittleEndian, &block)
		if block == 0 {
			break
		}
		obj.Pages = append(obj.Pages, block)
	}

	return obj
}

func (obj *Object) ReadData() []byte {
	data := make([]byte, 0, obj.Length)

	for _, page := range obj.Pages {
		pd, _ := obj.Database.ReadPage(uint(page))
		end := obj.Length - uint64(len(data))
		if end >= uint64(len(pd)) {
			end = uint64(len(pd))
		}

		data = append(data, pd[:end]...)
	}

	return data
}
