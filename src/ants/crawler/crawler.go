package crawler

import (
	"ants/util"
	"spiders"
)

type Crawler struct {
	SpiderList []*util.Spider
}

func (this *Crawler) LoadSpiders() {
	this.SpiderList = make([]*util.Spider, 1)
	deadLoopTest := spiders.MakeDeadLoopSpider()
	this.SpiderList = append(this.SpiderList, deadLoopTest)
}
