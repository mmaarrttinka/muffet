package app

import (
	"regexp"
	"time"
)

type FetcherOptions struct {
	Concurrency      int
	ExcludedPatterns []*regexp.Regexp
	Headers          map[string]string
	IgnoreFragments  bool
	FollowURLParams  bool
	MaxRedirections  int
	Timeout          time.Duration
	OnePageOnly      bool
}

func (o *FetcherOptions) Initialize() {
	if o.Concurrency <= 0 {
		o.Concurrency = DefaultConcurrency
	}

	if o.MaxRedirections <= 0 {
		o.MaxRedirections = DefaultMaxRedirections
	}

	if o.Timeout <= 0 {
		o.Timeout = DefaultTimeout
	}
}
