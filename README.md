### ants-go
open source, restful, distributed crawler engine
### comming up

* Persistence
* Dynamic

### design of ants-go
##### ants
I wrote a crawler engine named [ants](https://github.com/wcong/ants) in python base on [scrapy](https://github.com/scrapy/scrapy). But sometimes, dynamic language is chaos.
So I start to write it in a compile language. 
##### scrapy
I design the crawler framework  by imitating  [scrapy](https://github.com/scrapy/scrapy).
such as downloader,scraper,and the way user write customize spider,
but in a compile way
##### elasticsearch
I design my distributed architecture by imitating [elasticsearch](https://github.com/elasticsearch/elasticsearch).
it spire me to do a engine for distributed crawler
### requirement
``` shell
go get github.com/PuerkitoBio/goquery
go get github.com/go-sql-driver/mysql
```
### install

``` shell
go get github.com/wcong/ants-go
go install github.com/wcong/ants-go
```

### run

``` shell
cd bin
./ants-go
```

#### cluster in one computer
to test cluster in one computer,you can run it from different port in different terminal

one node,use the default port tcp 8300 http 8200

``` shell
cd bin
./ants-go
```

the other node set tcp port and http port

``` shell
cd bin
./ants-go -tcp 9300 -http 9200
```
#### flags
there are some flags you can set,check out the help message

``` shell
./ants-go -h
./ants-go -help
```

### Customize spider
1.	go to *spiders*
2.	write your spiders follow the example *deap_loop_spider.go* or go to the [spider page](./SPIDER.md)
3.	add you spider to spiderMap,follow the example in *LoadAllSpiders* in *load_all_spider.go*
4.	install again
