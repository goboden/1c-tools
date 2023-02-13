package v8repo

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

type Field struct {
	Type          string
	Empty         bool
	Length        uint16
	Accuracy      uint16
	CaseSensitive bool
	Size          uint16
	Offset        uint16
}

type Table struct {
	db           *Database
	Fields       map[string]Field
	FieldsIndex  []string
	RecordLength uint16
	DataPage     uint16
	IndexPage    uint16
	BlobPage     uint16
}

type TableReader struct {
	Data         []byte
	Rows         uint16
	Fields       map[string]Field
	FieldsIndex  []string
	RecordLength uint16
}

func NewTable(descr string) (string, *Table) {
	table := &Table{}
	table.Fields = make(map[string]Field)
	table.FieldsIndex = make([]string, 0)
	table.RecordLength = 1

	list, _ := NewListTree().Load(descr)

	name, _ := list.GetValue(0)
	name = strings.Trim(name, "\"")

	files, _ := list.GetValues(5)
	table.DataPage = StringToUint16(files[1])
	table.BlobPage = StringToUint16(files[2])
	table.IndexPage = StringToUint16(files[3])

	offset := uint16(1)

	fields, _ := list.Get(2)
	for i := 1; i < fields.Length(); i++ {
		fd, _ := fields.GetValues(i)

		field := &Field{Type: strings.Trim(fd[1], "\"")}

		fname := strings.Trim(fd[0], "\"")

		field.Length = StringToUint16(fd[3])
		field.Accuracy = StringToUint16(fd[4])

		if fd[2] == "1" {
			field.Empty = true
		}

		if fd[5] == "1" {
			field.CaseSensitive = true
		}

		field.Reacalculate(offset)
		offset = field.Offset + field.Size

		table.Fields[fname] = *field
		table.FieldsIndex = append(table.FieldsIndex, fname)
		table.RecordLength += field.Size
	}

	return name, table
}

func (f *Field) Reacalculate(offset uint16) {
	f.Offset = offset

	if f.Empty {
		f.Size += 1
	}

	switch f.Type {
	case "B": // Binary
		f.Size += f.Length
	case "L": // Bool
		f.Size += 1
	case "N": // Number
		f.Size += (f.Length + 2) / 2
	case "NC": // Fixed string
		f.Size += f.Length * 2
	case "NVC": // Variable string
		f.Size += f.Length*2 + 2
	case "RV": // Version
		f.Size += 16
	case "I": // Image
		f.Size += 8
	case "T": // Text
		f.Size += 8
	case "DT": // Datetime
		f.Size += 7
	case "NT": // Memo
		f.Size += 8
	default:
		panic(f.Type)
	}
}

func (tr *TableReader) ReadData(row uint16, fname string) []byte {
	field := tr.Fields[fname]
	begin := tr.RecordLength*row + field.Offset
	return tr.Data[begin : begin+field.Size]
}

func (tr *TableReader) ReadString(row uint16, fname string) string {
	field := tr.Fields[fname]
	begin := tr.RecordLength*row + field.Offset
	fdata := ReadFieldData(tr.Data[begin:begin+field.Size], field)
	return fdata
}

func (tr *TableReader) ReadUID(row uint16, fname string) string {
	field := tr.Fields[fname]
	fdata := tr.ReadString(row, fname)
	if field.Type == "B" && field.Length >= 16 {
		dec, _ := base64.StdEncoding.DecodeString(fdata)
		f1 := fmt.Sprintf("%x%x%x%x", dec[12], dec[13], dec[14], dec[15])
		f2 := fmt.Sprintf("%x%x", dec[10], dec[11])
		f3 := fmt.Sprintf("%x%x", dec[8], dec[9])
		f4 := fmt.Sprintf("%x%x", dec[0], dec[1])
		f5 := fmt.Sprintf("%x%x%x%x%x%x", dec[2], dec[3], dec[4], dec[5], dec[6], dec[7])
		fdata = fmt.Sprintf("%s-%s-%s-%s-%s", f1, f2, f3, f4, f5)
	}
	return fdata
}

func (tr *TableReader) ReadInteger(row uint16, fname string) uint16 {
	field := tr.Fields[fname]
	begin := tr.RecordLength*row + field.Offset
	fdata := ReadFieldData(tr.Data[begin:begin+field.Size], field)
	conv, _ := strconv.ParseUint(fdata, 10, 16)
	return uint16(conv)
}

func (tr *TableReader) FindString(fname string, find string) []uint16 {
	rows := make([]uint16, 0)

	for i := 1; i <= int(tr.Rows); i++ {
		fdata := tr.ReadString(uint16(i), fname)
		if fdata == find {
			rows = append(rows, uint16(i))
		}
	}

	return rows
}

func (tr *TableReader) FindUID(fname string, find string) []uint16 {
	rows := make([]uint16, 0)

	for i := 1; i <= int(tr.Rows); i++ {
		fdata := tr.ReadUID(uint16(i), fname)
		if fdata == find {
			rows = append(rows, uint16(i))
		}
	}

	return rows
}

func StringToUint16(s string) uint16 {
	vl, _ := strconv.ParseUint(s, 10, 16)
	return uint16(vl)
}

func ReadFieldData(data []byte, field Field) string {
	switch field.Type {
	case "B": // Binary
		return ReadB(data)
	case "L": // Bool
		return ReadL(data)
	case "N": // Number
		return ReadN(data, field.Length, field.Accuracy)
	case "NC": // Fixed string
		return ReadNC(data)
	case "NVC": // Variable string
		return ReadNVC(data)
	case "RV": // Version
		// f.Size += 16
	case "I": // Image
		// f.Size += 8
	case "T": // Text
		// f.Size += 8
	case "DT": // Datetime
		// f.Size += 7
	case "NT": // Memo
		// f.Size += 8
	default:
		panic(field.Type)
	}

	return ""
}

func ReadB(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func ReadL(data []byte) string {
	var value uint8
	binary.Read(bytes.NewReader(data[:1]), binary.LittleEndian, &value)
	if value == 1 {
		return "T"
	}
	return "F"
}

func ReadN(data []byte, length uint16, accuracy uint16) string {
	bcd := fmt.Sprintf("%x", data)
	return bcd[1 : length+1]
}

func ReadNVC(data []byte) string {
	var length uint16
	binary.Read(bytes.NewReader(data[:2]), binary.LittleEndian, &length)

	return ReadNC(data[2 : (length+1)*2])
}

func ReadNC(data []byte) string {
	var scode uint16
	result := ""
	for i := 0; i < len(data)/2; i++ {
		binary.Read(bytes.NewReader(data[i*2:i*2+2]), binary.LittleEndian, &scode)
		result += string(scode) // Ok
	}
	return result
}
