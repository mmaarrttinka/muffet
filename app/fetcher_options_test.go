package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetcherOptionsInitialize(t *testing.T) {
	o := FetcherOptions{}
	o.Initialize()

	assert.Equal(t, DefaultConcurrency, o.Concurrency)
	assert.Equal(t, DefaultMaxRedirections, o.MaxRedirections)
	assert.Equal(t, DefaultTimeout, o.Timeout)
}
