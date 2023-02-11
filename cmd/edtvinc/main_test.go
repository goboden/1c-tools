package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_maxver(t *testing.T) {
	type args struct {
		ver1 string
		ver2 string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"FirstMax", args{"0.0.0.1", "0.0.0.2"}, "0.0.0.2", false},
		{"SecondMax", args{"0.0.0.2", "0.0.0.1"}, "0.0.0.2", false},
		{"FirstErr", args{"0.0.0.s", "0.0.0.2"}, "", true},
		{"SecondErr", args{"0.0.0.1", "0.0.0.s"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := maxver(tt.args.ver1, tt.args.ver2)
			if (err != nil) != tt.wantErr {
				t.Errorf("maxver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("maxver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextVer(t *testing.T) {
	ver, err := nextver("0.0.0.1")
	assert.Equal(t, "0.0.0.2", ver)
	assert.Nil(t, err)

	_, err = nextver("Ver")
	assert.NotNil(t, err)
}

func TestNumVer(t *testing.T) {
	ver, err := numver("1.2.3.4")
	assert.Equal(t, [4]int{1, 2, 3, 4}, ver)
	assert.Nil(t, err)

	_, err = numver("1.2.3.Ver")
	assert.NotNil(t, err)

	_, err = numver("Ver")
	assert.NotNil(t, err)

}
