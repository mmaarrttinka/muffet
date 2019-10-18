package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDonePageSet(t *testing.T) {
	newDonePageSet()
}

func TestDonePageSetAdd(t *testing.T) {
	s := newDonePageSet()
	assert.False(t, s.Add("foo"))
	assert.True(t, s.Add("foo"))
}
