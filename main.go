package main

import (
	"flag"
	"github.com/wcong/ants-go/ants/action"
	AHttp "github.com/wcong/ants-go/ants/action/http"
	"github.com/wcong/ants-go/ants/action/rpc"
	"github.com/wcong/ants-go/ants/action/watcher"
	"github.com/wcong/ants-go/ants/crawler"
	"github.com/wcong/ants-go/ants/http"
	"github.com/wcong/ants-go/ants/node"
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
func initCluster(settings *util.Settings, rpcClient action.RpcClientAnts, node *node.Node) {
	node.Join()
	isClusterExist := false
	if len(settings.NodeList) > 0 {
		for _, nodeInfo := range settings.NodeList {
			nodeSettings := strings.Split(nodeInfo, ":")
			ip := nodeSettings[0]
			port, _ := strconv.Atoi(nodeSettings[1])
			if ip == node.NodeInfo.Ip && port == node.NodeInfo.Port {
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
		node.MakeMasterNode(node.NodeInfo.Name)
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
	Node := node.NewNode(setting, resultQuene)
	var rpcClient action.RpcClientAnts = rpc.NewRpcClient(Node)
	var distributer action.Watcher = watcher.NewDistributer(Node, rpcClient)
	var reporter action.Watcher = watcher.NewReporter(Node, rpcClient, resultQuene, distributer)
	rpc.NewRpcServer(Node, setting.TcpPort, rpcClient, reporter, distributer)
	router := AHttp.NewRouter(Node, reporter, distributer, rpcClient)
	httpServer := http.NewHttpServer(setting, router)
	httpServer.Start(wg)
	initCluster(setting, rpcClient, Node)
	rpcClient.Start()
	wg.Wait()
}
