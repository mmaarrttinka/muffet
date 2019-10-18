package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"os"

	"github.com/valyala/fasthttp"
)

func main() {
	s, err := command(os.Args[1:], os.Stdout)

	if err != nil {
		fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(s)
}

func command(ss []string, w io.Writer) (int, error) {
	args, err := getArguments(ss)

	if err != nil {
		return 0, err
	}

	h := &fasthttp.Client{
		MaxConnsPerHost: args.Concurrency,
		ReadBufferSize:  args.BufferSize,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: args.SkipTLSVerification,
		},
	}

	c, err := newChecker(args.URL, h, checkerOptions{
		fetcherOptions{
			args.Concurrency,
			args.ExcludedPatterns,
			args.Headers,
			args.IgnoreFragments,
			args.FollowURLParams,
			args.MaxRedirections,
			args.Timeout,
			args.OnePageOnly,
		},
		args.FollowRobotsTxt,
		args.FollowSitemapXML,
		args.FollowURLParams,
	})

	if err != nil {
		return 0, err
	}

	go c.Check()

	s := 0

	for r := range c.Results() {
		if !r.OK() || args.Verbose {
			fprintln(w, r.String(args.Verbose))
		}

		if !r.OK() {
			s = 1
		}
	}

	return s, nil
}

func fprintln(w io.Writer, xs ...interface{}) {
	if _, err := fmt.Fprintln(w, xs...); err != nil {
		panic(err)
	}
}
