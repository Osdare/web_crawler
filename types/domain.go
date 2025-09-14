package types

type Domain struct {
	Name        string
	CrawlDelay  int
	LastCrawled int
	Disallowed  bool
}
