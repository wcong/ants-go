package http

import (
	"github.com/wcong/ants-go/ants/node"
)

// welcome struct
type WelcomeInfo struct {
	Message  string
	Greeting string
	Time     string
}

// result of start spider
type StartSpiderResult struct {
	Success    bool
	Message    string
	Spider     string
	Time       string
	MasterNode *node.NodeInfo
}
