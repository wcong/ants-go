package node

import (
	"time"
)

/*
*	record spider job status
*
**/

type SpiderStatus struct {
	StartTime  time.Time
	EndTime    time.Time
	Name       string
	RequestNum int
}

type JobStatus struct {
	FinishedJobs []*SpiderStatus
	RunningJobs  map[string]*SpiderStatus
}
