package app

import (
	"errors"
	"sync"

	"github.com/fatih/color"
	"github.com/valyala/fasthttp"
)

// Checker checks linkes in pages.
type Checker struct {
	fetcher
	daemons      daemons
	urlInspector urlInspector
	results      chan pageResult
	donePages    donePageSet
}

// NewChecker creates a new link checker.
func NewChecker(s string, c *fasthttp.Client, o CheckerOptions) (Checker, error) {
	o.Initialize()

	f := newFetcher(c, o.FetcherOptions)
	r, err := f.Fetch(s)

	if err != nil {
		return Checker{}, err
	}

	p, ok := r.Page()

	if !ok {
		return Checker{}, errors.New("non-HTML page")
	}

	ui, err := newURLInspector(c, p.URL().String(), o.FollowRobotsTxt, o.FollowSitemapXML)

	if err != nil {
		return Checker{}, err
	}

	ch := Checker{
		f,
		newDaemons(o.Concurrency),
		ui,
		make(chan pageResult, o.Concurrency),
		newDonePageSet(),
	}

	ch.addPage(p)

	return ch, nil
}

func (c Checker) Results() <-chan pageResult {
	return c.results
}

func (c Checker) Check() {
	c.daemons.Run()

	close(c.results)
}

func (c Checker) checkPage(p *page) {
	us := p.Links()

	sc := make(chan string, len(us))
	ec := make(chan string, len(us))
	w := sync.WaitGroup{}

	for u, err := range us {
		if err != nil {
			ec <- formatLinkError(u, err)
			continue
		}

		w.Add(1)

		go func(u string) {
			defer w.Done()

			r, err := c.fetcher.Fetch(u)

			if err == nil {
				sc <- formatLinkSuccess(u, r.StatusCode())
			} else {
				ec <- formatLinkError(u, err)
			}

			// only consider adding the page to the list if we're recursing
			if !c.fetcher.options.OnePageOnly {
				if p, ok := r.Page(); ok && c.urlInspector.Inspect(p.URL()) {
					c.addPage(p)
				}
			}
		}(u)
	}

	w.Wait()

	c.results <- newPageResult(p.URL().String(), stringChannelToSlice(sc), stringChannelToSlice(ec))
}

func (c Checker) addPage(p *page) {
	if !c.donePages.Add(p.URL().String()) {
		c.daemons.Add(func() { c.checkPage(p) })
	}
}

func stringChannelToSlice(sc <-chan string) []string {
	ss := make([]string, 0, len(sc))

	for i := 0; i < cap(ss); i++ {
		ss = append(ss, <-sc)
	}

	return ss
}

func formatLinkSuccess(u string, s int) string {
	return color.GreenString("%v", s) + "\t" + u
}

func formatLinkError(u string, err error) string {
	return color.RedString(err.Error()) + "\t" + u
}
