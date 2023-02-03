package v8unpack

import (
	"fmt"
	"os"
	"testing"
)

func TestContainer_ReadFile(t *testing.T) {
	input := "./testdata/ExtForm.epf"
	file, err := os.Open(input)
	if err != nil {
		t.Fatal()
	}
	defer file.Close()

	reader := NewFileReader(file)
	rcnt := ReadRootContainer(reader)

	fmt.Printf("%v\n", rcnt.index)

	fl, err := rcnt.ReadFile("somename", false)
	if err == nil {
		fmt.Println(fl.IsContainer)
	}
}
