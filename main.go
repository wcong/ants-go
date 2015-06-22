package main

import (
	"flag"
	"github.com/wcong/ants-go/ants/action"
	AHttp "github.com/wcong/ants-go/ants/action/http"
	"github.com/wcong/ants-go/ants/action/rpc"
	"github.com/wcong/ants-go/ants/action/watcher"
	clusterSupport "github.com/wcong/ants-go/ants/cluster/support"
	"github.com/wcong/ants-go/ants/crawler"
	"github.com/wcong/ants-go/ants/http"
	"github.com/wcong/ants-go/ants/node"
	nodeSupport "github.com/wcong/ants-go/ants/node/support"
	"github.com/wcong/ants-go/ants/util"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

// init command line param
func initFlag(settings *util.Settings) {
	flag.IntVar(&settings.TcpPort, "tcp", settings.TcpPort, "tcp port")
	flag.IntVar(&settings.HttpPort, "http", settings.HttpPort, "http port")
	flag.IntVar(&settings.DownloadInterval, "di", settings.DownloadInterval, "Download Interval in second")
	flag.StringVar(&settings.ConfigFile, "c", settings.ConfigFile, "config file path")
}

// load setting from json file and input
func MakeSettings() *util.Settings {
	settings := util.NewSettings()
	initFlag(settings)
	flag.Parse()
	if settings.ConfigFile != "" {
		util.LoadSettingFromFile(settings.ConfigFile, settings)
	}
	return settings
}

// try to join the cluster,
// if there is no cluster,make itself master
func initCluster(settings *util.Settings, rpcClient action.RpcClientAnts, node node.Node) {
	node.Join()
	isClusterExist := false
	if len(settings.NodeList) > 0 {
		for _, nodeInfo := range settings.NodeList {
			nodeSettings := strings.Split(nodeInfo, ":")
			ip := nodeSettings[0]
			port, _ := strconv.Atoi(nodeSettings[1])
			if ip == node.GetNodeInfo().Ip && port == node.GetNodeInfo().Port {
				continue
			}
			err := rpcClient.LetMeIn(ip, port)
			if err == nil {
				isClusterExist = true
			}
		}
	} else {

	}
	if !isClusterExist {
		node.MakeMasterNode(node.GetNodeInfo().Name)
	}
	node.Ready()
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("let us go shipping")
	var wg sync.WaitGroup
	wg.Add(1)
	setting := MakeSettings()
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "-help") {
		flag.PrintDefaults()
		return
	}
	resultQuene := crawler.NewResultQuene()
	defaultNode := nodeSupport.NewDefaultNode(setting, resultQuene)
	defaultCluster := clusterSupport.NewDefaultCluster(setting)
	defaultNode.Init(defaultCluster)
	defaultCluster.Init(defaultNode.GetNodeInfo())
	var rpcClient action.RpcClientAnts = rpc.NewRpcClient(defaultNode, defaultCluster)
	var distributer action.Watcher = watcher.NewDistributer(defaultNode, defaultCluster, rpcClient)
	var reporter action.Watcher = watcher.NewReporter(defaultNode, defaultCluster, rpcClient, resultQuene, distributer)
	rpc.NewRpcServer(defaultNode, defaultCluster, setting.TcpPort, rpcClient, reporter, distributer)
	router := AHttp.NewRouter(defaultNode, defaultCluster, reporter, distributer, rpcClient)
	httpServer := http.NewHttpServer(setting, router)
	httpServer.Start(wg)
	initCluster(setting, rpcClient, defaultNode)
	rpcClient.Start()
	wg.Wait()
}
