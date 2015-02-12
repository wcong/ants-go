package crawler

import (
	"time"
)

type SpiderStatus struct {
	Name      string
	Crawled   int
	Running   int
	Waiting   int
	StartTime time.Time
	EndTime   time.Time
}

type CrawlerStatus struct {
	CrawledSpider []*SpiderStatus
	RuningSpider  []*SpiderStatus
}
