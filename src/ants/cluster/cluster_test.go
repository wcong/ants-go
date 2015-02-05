package cluster

import (
	"testing"
)

func TestCluster(t *testing.T) {
	testCluster := Cluster{"test", nil, nil}
	testCluster.SetName("test")
}
