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

func main() {
	var input, output string

	flag.Parse()

	input = flag.Arg(0)
	output = flag.Arg(1)

	if input == "" || output == "" {
		fmt.Println("Usage: modunpack [Input file] [Output dir]")
		os.Exit(1)
	}

	if info, err := os.Stat(input); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("File '%s' not exist\n", input)
		os.Exit(1)
	} else {
		if info.IsDir() {
			fmt.Printf("'%s' is a directory\n", input)
			os.Exit(1)
		}
	}

	// if info, err := os.Stat(output); errors.Is(err, os.ErrNotExist) {
	// 	fmt.Printf("Directory '%s' not exist\n", output)
	// 	os.Exit(1)
	// } else {
	// 	if !info.IsDir() {
	// 		fmt.Printf("'%s' is not a directory\n", output)
	// 		os.Exit(1)
	// 	}
	// }

	unpack(input, output)
}

func unpack(input string, output string) {
	extname := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))

	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	reader := v8unpacker.NewFileReader(file)
	root := v8unpacker.ReadRootContainer(reader)

	modules, err := v8unpacker.FindModules(root)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	for key, value := range modules {
		Save(output+"/"+extname+"/", key+".bsl", value)
	}

}

func Save(path string, filename string, content string) {
	os.MkdirAll(path, 0666)
	err := os.WriteFile(path+filename, []byte(content), 0666)
	if err != nil {
		fmt.Println(err.Error(), filename)
	}
}
