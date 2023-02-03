package v8repo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

type Database struct {
	Filename   string
	Signature  string // 8
	Version    string // 1+1+1+1
	Length     uint32 // 4
	U          int32  // 4
	PageLength uint32 // 4
	Tables     map[string]*Table
}

func NewDatabase(filename string) (*Database, error) {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	db := &Database{Filename: filename}
	db.Tables = make(map[string]*Table)

	db.ReadHeader()
	db.ReadTables()

	return db, nil
}

func (db *Database) ReadData(begin uint, length uint) ([]byte, error) {
	file, err := os.Open(db.Filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := NewFileReader(file)
	data := reader.ReadFragment(begin, length)

	return data, nil
}

func (db *Database) ReadPage(number uint) ([]byte, error) {
	bl := uint(db.PageLength)
	begin := number * bl
	// fmt.Printf("Block %d: %d / %x - %d / %x\n", number, begin, begin, begin+bl, begin+bl)
	block, err := db.ReadData(begin, bl)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (db *Database) ReadHeader() {
	data, _ := db.ReadData(0, 32)

	db.Signature = string(data[0:8])

	version := data[8:12]
	db.Version = fmt.Sprintf("%d.%d.%d.%d", version[0], version[1], version[2], version[3])

	binary.Read(bytes.NewReader(data[12:16]), binary.LittleEndian, &db.Length)
	binary.Read(bytes.NewReader(data[16:20]), binary.LittleEndian, &db.U)
	binary.Read(bytes.NewReader(data[20:24]), binary.LittleEndian, &db.PageLength)
}

func (db *Database) ReadTables() {
	tables := db.ReadBlob(3)
	for i := 0; i < len(tables.Data); i++ {
		name, table := NewTable(string(tables.Data[i]))
		table.db = db
		db.Tables[name] = table
	}

	// for k, v := range db.Tables {
	// 	fmt.Printf("Table: %s (%d, %d, %d) L=%d\n", k, v.DataPage, v.IndexPage, v.BlobPage, v.RecordLength)
	// 	for m, o := range v.Fields {
	// 		fmt.Printf("\tField: %10s: %+v\n", m, o)
	// 	}
	// }

}

func (db *Database) NewTableReader(tname string) *TableReader {
	table := db.Tables[tname]

	fmt.Printf("Table: %s L=%d\n", tname, table.RecordLength)
	for i, fname := range table.FieldsIndex {
		field := table.Fields[fname]
		fmt.Printf(".%2d %10s (%3s): +%3d %+v\n", i+1, fname, field.Type, field.Offset, field)
	}

	obj := NewObject(db, uint32(table.DataPage))
	data := obj.ReadData()

	rows := uint16(len(data))/table.RecordLength - 1

	return &TableReader{Data: data, Fields: table.Fields, FieldsIndex: table.FieldsIndex, RecordLength: table.RecordLength, Rows: rows}
}

func (db *Database) ReadBlob(page uint) *Blob {
	data, _ := db.ReadPage(page)
	blob := NewBlob(data)

	// 0 Free blocks

	// 1 Header

	return blob
}

func (db *Database) ReadTableData(tname string) {
	reader := db.NewTableReader(tname)

	nr := int(reader.Rows)
	// nr = 3 // !!!

	for i := 1; i <= nr; i++ {
		fmt.Printf("___r[%2d]___\n", i)
		for fi, fname := range reader.FieldsIndex {
			fdata := reader.ReadString(uint16(i), fname)
			fmt.Printf("%2d %10s [%s]\n", fi, fname, fdata)
		}
	}

}
