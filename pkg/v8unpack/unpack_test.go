package v8unpack

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileReader_ReadFragment(t *testing.T) {
	input := "./testdata/data.hex"
	file, err := os.Open(input)
	if err != nil {
		t.Fatal()
	}
	defer file.Close()

	var f []byte
	var ok bool
	reader := NewFileReader(file)

	f, ok = reader.ReadFragment(0, 250)
	assert.Equal(t, 250, len(f))
	assert.Equal(t, true, ok)

	f, ok = reader.ReadFragment(256, 8)
	assert.Equal(t, 0, len(f))
	assert.Equal(t, false, ok)

	f, ok = reader.ReadFragment(250, 8)
	assert.Equal(t, true, ok)
	assert.Equal(t, []byte{251, 252, 253, 254, 255, 0}, f)
}

func TestBytesReader_ReadFragment(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}

	var f []byte
	var ok bool
	reader := NewBytesReader(data)

	f, ok = reader.ReadFragment(0, 8)
	assert.Equal(t, 8, len(f))
	assert.Equal(t, true, ok)

	f, ok = reader.ReadFragment(9, 16)
	assert.Equal(t, 0, len(f))
	assert.Equal(t, false, ok)

	f, ok = reader.ReadFragment(5, 3)
	assert.Equal(t, []byte{6, 7, 8}, f)
	assert.Equal(t, true, ok)

	f, ok = reader.ReadFragment(5, 16)
	assert.Equal(t, []byte{6, 7, 8}, f)
	assert.Equal(t, true, ok)
}

func Test_bytesToString(t *testing.T) {
	a := bytesToString([]byte{
		0x37, 0x00, 0x66, 0x00, 0x38, 0x00, 0x38, 0x00,
		0x30, 0x00, 0x39, 0x00, 0x30, 0x00, 0x31, 0x00,
		0x2d, 0x00, 0x36, 0x00, 0x32, 0x00, 0x38, 0x00,
		0x66, 0x00, 0x2d, 0x00, 0x34, 0x00, 0x63, 0x00,
		0x32, 0x00, 0x36, 0x00, 0x2d, 0x00, 0x62, 0x00,
		0x38, 0x00, 0x62, 0x00, 0x36, 0x00, 0x2d, 0x00,
		0x33, 0x00, 0x35, 0x00, 0x31, 0x00, 0x63, 0x00,
		0x35, 0x00, 0x37, 0x00, 0x38, 0x00, 0x38, 0x00,
		0x37, 0x00, 0x39, 0x00, 0x61, 0x00, 0x62, 0x00,
	})
	assert.Equal(t, "7f880901-628f-4c26-b8b6-351c578879ab", a)
}

func Test_bytesToAdr(t *testing.T) {
	a := bytesToAdr([]byte{48, 48, 48, 48, 48, 50, 48, 48})
	assert.Equal(t, v8address(0x200), a)
}
