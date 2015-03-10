package node

import (
	"ants/crawler"
	"ants/http"
	"log"
)

/*
*	this is rpc method of request
**/
func (this *RPCer) distribute(nodeName string, request *http.Request) error {
	distributeRequest := &DistributeRequest{}
	distributeRequest.NodeInfo = this.node.NodeInfo
	distributeRequest.Request = request
	distributeReqponse := &DistributeReqponse{}
	err := this.connMap[nodeName].Call("RPCer.AcceptRequest", distributeRequest, distributeReqponse)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(distributeReqponse.Result)
	}
	return err
}

// expose method ,for accept method
func (this *RPCer) AcceptRequest(request *DistributeRequest, response *DistributeReqponse) error {
	this.node.AcceptRequest(request.Request)
	return nil
}

// for slave send crawl result to master
func (this *RPCer) reportResult(nodeName string, result *crawler.ScrapeResult) error {
	reportRequest := &ReportRequest{}
	reportRequest.NodeInfo = this.node.NodeInfo
	reportRequest.ScrapeResult = result
	reportResponse := &ReportResponse{}
	err := this.connMap[nodeName].Call("RPCer.AcceptResult", reportRequest, reportResponse)
	if err != nil {
		log.Println(err)
	}
	return err
}

// for master accept crawl result
func (this *RPCer) AcceptResult(request *ReportRequest, response *ReportResponse) error {
	this.node.AcceptResult(request.ScrapeResult)
	return nil
}
