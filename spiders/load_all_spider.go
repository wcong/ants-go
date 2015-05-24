package spiders

import (
	"github.com/wcong/ants-go/ants/spiders"
)

func LoadAllSpiders() map[string]*spiders.Spider {
	spiderMap := make(map[string]*spiders.Spider)
	deadLoopTest := MakeDeadLoopSpider()
	spiderMap[deadLoopTest.Name] = deadLoopTest
	dumpTestSpider := MakeDumpTestSpider()
	spiderMap[dumpTestSpider.Name] = dumpTestSpider
	zhihuSpider := MakeZhiHuSpider()
	spiderMap[zhihuSpider.Name] = zhihuSpider
	muiltiplySpiders := MakeMuiltiplySpiders()
	spiderMap[muiltiplySpiders.Name] = muiltiplySpiders
	stockSpider := MakeStockSpider()
	spiderMap[stockSpider.Name] = stockSpider
	stockPriceSpider := MakeStockPriceSpider()
	spiderMap[stockPriceSpider.Name] = stockPriceSpider
	return spiderMap
}
