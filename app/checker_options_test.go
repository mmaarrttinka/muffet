package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckerOptionsInitialize(t *testing.T) {
	o := CheckerOptions{}
	o.Initialize()

	assert.Equal(t, DefaultConcurrency, o.Concurrency)
}
