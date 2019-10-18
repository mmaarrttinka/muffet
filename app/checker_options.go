package app

type CheckerOptions struct {
	FetcherOptions
	FollowRobotsTxt,
	FollowSitemapXML,
	FollowURLParams bool
}
