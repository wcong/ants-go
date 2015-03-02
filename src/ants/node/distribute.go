package node

import (
	"ants/http"
)

/*
what a Distributer do
*	status,running|parse|stop
*	distribute a request,by some strategy
*	*
*/
const (
	DISTRIBUTE_RUNING = iota
	DISTRIBUTE_PAUSE
	DISTRIBUTE_STOP
)

type Distributer struct {
	Status    int
	Cluster   *Cluster
	LastIndex int
}

func NewDistributer(cluster *Cluster) *Distributer {
	return &Distributer{DISTRIBUTE_STOP, cluster, 0}
}

func (this *Distributer) IsStop() bool {
	return this.Status == DISTRIBUTE_STOP
}

func (this *Distributer) IsPause() bool {
	return this.Status == DISTRIBUTE_PAUSE
}
func (this *Distributer) Run() {
	this.Status = DISTRIBUTE_RUNING
}
func (this *Distributer) Pause() {
	this.Status = DISTRIBUTE_PAUSE
}
func (this *Distributer) Stop() {
	this.Status = DISTRIBUTE_STOP
}

// if cookiejar > 0 means it require cookie context ,so we should send it to where it come from
// else distribute it by order
func (this *Distributer) Distribute(request *http.Request) string {
	if request.CookieJar > 0 {
		return request.NodeName
	} else {
		if this.LastIndex > len(this.Cluster.ClusterInfo.NodeList) {
			this.LastIndex = 0
		}
		nodeName := this.Cluster.ClusterInfo.NodeList[this.LastIndex].Name
		request.NodeName = nodeName
		this.LastIndex += 1
		return nodeName
	}
}
