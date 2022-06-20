package mpc

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestGroupTest(t *testing.T) {
	GroupTest()
	assert.Equal(t, 8, 8)
	test2()
}
