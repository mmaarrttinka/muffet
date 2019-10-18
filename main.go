package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"os"

	"github.com/raviqqe/muffet/app"
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

	c, err := app.NewChecker(args.URL, h, app.CheckerOptions{
		FetcherOptions: app.FetcherOptions{
			Concurrency:      args.Concurrency,
			ExcludedPatterns: args.ExcludedPatterns,
			Headers:          args.Headers,
			IgnoreFragments:  args.IgnoreFragments,
			FollowURLParams:  args.FollowURLParams,
			MaxRedirections:  args.MaxRedirections,
			Timeout:          args.Timeout,
			OnePageOnly:      args.OnePageOnly,
		},
		FollowRobotsTxt:  args.FollowRobotsTxt,
		FollowSitemapXML: args.FollowSitemapXML,
		FollowURLParams:  args.FollowURLParams,
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
