package action

import (
	"net/http"
)

/*
*	http server for ants should server those function
**/

// http server
type HttpServer interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Welcome(w http.ResponseWriter, r *http.Request)
}

// http server for cluster
type HttpServerCluster interface {
	Cluster(w http.ResponseWriter, r *http.Request)
}

// http server for crawler
type HttpServerCrawler interface {
	Spiders(w http.ResponseWriter, r *http.Request)
	Crawl(w http.ResponseWriter, r *http.Request)
	CrawlCluster(w http.ResponseWriter, r *http.Request)
	CrawlNode(w http.ResponseWriter, r *http.Request)
}

// ants http server
type HttpServerAnts interface {
	HttpServer
	HttpServerCluster
	HttpServerCrawler
}
