package types

type Domain struct {
	Name        string
	CrawlDelay  int
	LastCrawled int
	Allowed     []string
	Disallowed  []string
}
