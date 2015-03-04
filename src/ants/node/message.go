package node

import (
	"ants/http"
)

const (
	HADNLER_JOIN_REQUEST = iota
	HADNLER_JOIN_RESPONSE
	HADNLER_JOIN_EXAM // all node fix node list and master node
	HANDLER_SEND_MASTER_REQUEST
	HANDLER_SEND_REQUEST
	HANDLER_SEND_REQUEST_RESULT
	HANDLER_STOP_NODE
)

// what a message would do
// *		message type
// *		the request for crawl
// *		the result of crawl
// *		scraped request result
// *		the node which scrape it
type RequestMessage struct {
	Type            int
	Request         *http.Request
	CrawlResult     string // if success just empty string,or error reason
	ScrapedRequests []*http.Request
	NodeInfo        *NodeInfo
	ClusterInfo     *ClusterInfo
}
