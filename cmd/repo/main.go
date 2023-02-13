package main

import (
	"bytes"
	"compress/flate"
	"fmt"
	"os"

	"gitlab.helix.ru/terentev.d/git_tools/pkg/v8repo"
	"gitlab.helix.ru/terentev.d/git_tools/pkg/v8unpack"
)

// https://infostart.ru/1c/articles/19734/
// https://infostart.ru/1c/articles/536343/

func main() {
	rd := "./repo"
	fmt.Printf("Repository:  %s\n", rd)

	rep, err := NewRepository(rd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db := rep.Database

	// DatabaseInfo(db)

	// page, _ := db.ReadPage(3)
	// fmt.Printf("[%d]", page[0:256])

	// blob := db.ReadBlob(3)
	// fmt.Printf("[%d]", len(blob.Data))
	// for _, v := range blob.Data {
	// 	fmt.Printf("___\n[%s]\n", string(v))
	// }

	// list, _ := NewListTree().Load(string(blob.Data[1]))
	// fmt.Printf("%s\n", list.ToString())

	// root, _ := db.ReadPage(5)
	// obj := NewObject(root)
	// fmt.Printf("{%+v}", obj)

	// db.ReadTableData("USERS")
	// db.ReadTableData("OBJECTS")

	// reader := db.NewTableReader("USERS")

	// for i := 1; i <= int(reader.Rows); i++ {
	// 	fmt.Printf("User: %15s uid=%s\n", reader.ReadString(uint16(i), "NAME"), reader.ReadUID(uint16(i), "USERID"))
	// }

	// idx := reader.FindString("NAME", "Decompose")
	// for _, row := range idx {
	// 	fmt.Printf("FOUND User: %15s uid=%s\n", reader.ReadString(row, "NAME"), reader.ReadUID(row, "USERID"))
	// }

	// versions := db.NewTableReader("VERSIONS")
	// users := db.NewTableReader("USERS")

	// for i := uint16(1); i <= versions.Rows; i++ {
	// 	// uid := versions.Row(i).String("USERID")

	// 	uid := versions.ReadUID(i, "USERID")
	// 	urows := users.FindUID("USERID", uid)
	// 	fmt.Printf("%3d: %5d -> %s\n", i, versions.ReadInteger(i, "VERNUM"), users.ReadString(urows[0], "NAME"))
	// }

	// !!! EXTERNALS/EXTVERID = data/objects/[0]/[1:]
	// 9e ce c0 be 4e . c0 2c 21 0e 70 . 5b b6 6c ca f8 . 31 27 47 ca ++ ad -19+1
	// 01 ab cc 37 73 . 03 32 4c b4 1b . 1b 62 fd 6b d1 . 0e 20 79 c3 d9 ad -21

	// 15 af 5e 8f 18 . ae a8 39 7a e2 . 87 ef 2f 0f 23 . 7e ca 7b a4 ++ 10
	// 15 af 5e 8f 18 . ae a8 39 7a e2 . 87 ef 2f 0f 23 7e ca 7b a4 10

	ind := uint16(1)
	history := db.NewTableReader("EXTERNALS")
	ds := history.ReadString(ind, "EXTNAME")
	dt := history.ReadData(ind, "EXTVERID")
	dp := history.ReadString(ind, "DATAPACKED")
	fmt.Printf("-> %s [%s] %2x / %2x %d\n", dp, ds, dt[0], dt[1:], dt)

	// fpath := fmt.Sprintf("%s/data/objects/%x/%x", rep.Directory, dt[0], dt[1:])
	// fpath := fmt.Sprintf("%s/data/objects/15/2b70f955c1d6d387690b1c3a09404dc1c27111", rep.Directory)
	// fpath := fmt.Sprintf("%s/data/objects/89/a23f7ce3915d36abe0800d4c17ba6f784ec2f0", rep.Directory)
	fpath := fmt.Sprintf("%s/data/objects/6e/1cd76bc4c6bcf97729ab9b284639b1be9b40a4", rep.Directory)

	// file, err := os.Open(fpath)
	// if err != nil {
	// 	fmt.Println(fpath, err.Error())

	// 	os.Exit(1)
	// }
	// defer file.Close()

	// fr := v8unpacker.NewFileReader(file)
	// rt := v8unpacker.ReadContainer(fr)
	// for name := range rt.GetIndex() {
	// 	println(name)
	// }

	fdd, _ := os.ReadFile(fpath)
	fdef := deflate(fdd)
	// fmt.Printf("%s\n", fdef)

	br := v8unpacker.NewBytesReader(fdef)
	con := v8unpacker.ReadContainer(br)
	for name := range con.GetIndex() {
		println(name)
	}

}

func DatabaseInfo(db *Database) {
	fmt.Printf("Database:    %s\n", db.Filename)
	fmt.Printf("Signature:   %s\n", db.Signature)
	fmt.Printf("Version:     %s\n", db.Version)
	fmt.Printf("Length:      %d\n", db.Length)
	fmt.Printf("U:           %d\n", db.U)
	fmt.Printf("BlockLength: %d\n", db.PageLength)
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
