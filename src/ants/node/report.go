package node

import (
	"ants/crawler"
	"time"
)

/*
* report request result of scrapy
*	*	container list of report
*	*	dead loop for report result to master
**/

const (
	REPORT_STATUS_RUNNING = iota
	REPORT_STATUS_PAUSE
	REPORT_STATUS_STOP
)

type Reporter struct {
	Status      int
	ResultQuene *crawler.ResultQuene
	Node        *Node
}

func NewReporter(node *Node) *Reporter {
	resultQuene := crawler.NewResultQuene()
	return &Reporter{REPORT_STATUS_STOP, resultQuene, node}
}

func (this *Reporter) Start() {
	if this.Status == REPORT_STATUS_RUNNING {
		return
	}
	this.Status = REPORT_STATUS_RUNNING
	go this.Run()
}

func (this *Reporter) Pause() {
	this.Status = REPORT_STATUS_PAUSE
}

func (this *Reporter) Stop() {
	this.Status = REPORT_STATUS_STOP
}

func (this *Reporter) IsStop() bool {
	return this.Status == REPORT_STATUS_STOP
}

func (this *Reporter) Run() {
	for {
		if this.Status == REPORT_STATUS_PAUSE {
			time.Sleep(1 * time.Second)
			continue
		}
		if this.Status != REPORT_STATUS_RUNNING {
			break
		}
		result := this.ResultQuene.Pop()
		if result != nil {
			this.Node.ReportToMaster(result)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
