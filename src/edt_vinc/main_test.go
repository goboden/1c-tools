package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxVer(t *testing.T) {
	ver, err := maxver("0.0.0.1", "0.0.0.2")
	assert.Equal(t, "0.0.0.2", ver, "Wrong max version: "+ver)
	assert.Nil(t, err)

	ver, _ = maxver("0.0.0.2", "0.0.0.1")
	assert.Equal(t, "0.0.0.2", ver, "Wrong max version: "+ver)

	_, err = maxver("0.0.0.1s", "0.0.0.2")
	assert.NotNil(t, err)

	_, err = maxver("0.0.0.1", "0.0.0.2s")
	assert.NotNil(t, err)
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
