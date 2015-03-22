### ants-go 
open source, restful, distributed crawler engine  
### why
I wrote a crawler engine named [ants](https://github.com/wcong/ants) in python base on [scrapy](https://github.com/scrapy/scrapy). But sometimes, dynamic language is chaos.
So I start to write it in a compile language.  
### alpha
just use it for fun by now.
### requirement
``` shell
go get github.com/PuerkitoBio/goquery
go get github.com/go-sql-driver/mysql
```
### install

``` shell
go install src/ants/ants/ants.go
```

### run

``` shell
cd bin
./ants
```

there are some flags you can set,check out the help message

``` shell
./ants -h
./ants -help
```

### Customize spider
1.	go to *src/spiders*
2.	write your spiders follow the example *deap_loop_spider.go* or go to the [spider page](./SPIDER.md)
3.	add you spider to spiderMap,follow the example in *LoadAllSpiders* in *load_all_spider.go*
4.	install again

