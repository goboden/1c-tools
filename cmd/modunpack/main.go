package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"log"
	"path/filepath"
	"strings"

	"gitlab.helix.ru/terentev.d/git_tools/pkg/v8unpack"
)

type opts struct {
	input  string
	output string
}

func main() {
	opts := &opts{}
	opts.parse()

	if info, err := os.Stat(opts.input); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("File '%s' not exist\n", opts.input)
		os.Exit(1)
	} else {
		if info.IsDir() {
			fmt.Printf("'%s' is a directory\n", opts.input)
			os.Exit(1)
		}
	}

	unpack(opts.input, opts.output)
}

func unpack(input string, output string) {
	extname := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))

	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	reader := v8unpack.NewFileReader(file)
	root := v8unpack.ReadRootContainer(reader)

	modules, err := v8unpack.FindModules(root)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	for key, value := range modules {
		save(output+"/"+extname+"/", key+".bsl", value)
	}

}

func save(path string, filename string, content string) {
	os.MkdirAll(path, 0666)
	err := os.WriteFile(path+filename, []byte(content), 0666)
	if err != nil {
		fmt.Println(err.Error(), filename)
	}
}

func (o *opts) parse() {
	flag.Parse()

	o.input = flag.Arg(0)
	o.output = flag.Arg(1)

	if o.input == "" || o.output == "" {
		fmt.Println("Usage: modunpack [Input file] [Output dir]")
		os.Exit(1)
	}
}
