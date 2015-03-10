package node

import (
	"ants/util"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
)

const (
	RPC_TYPE = "tcp"
)

type RPCer struct {
	node     *Node
	settings *util.Settings
	connMap  map[string]*rpc.Client
}

func NewRPCer(node *Node, settings *util.Settings) *RPCer {
	connMap := make(map[string]*rpc.Client)
	return &RPCer{node, settings, connMap}
}

// start a rpc server
func (this *RPCer) server() {
	rpc.Register(this)
	listener, e := net.Listen(RPC_TYPE, ":"+strconv.Itoa(this.settings.TcpPort))
	if e != nil {
		log.Println(e)
		return
	}
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Println(err)
		} else {
			log.Println("new connection")
			go jsonrpc.ServeConn(conn)
		}
	}
}

func (this *RPCer) start() {
	go this.server()
}

// start a rpc client
// if ok , store if
func (this *RPCer) dial(ip string, port int) (*rpc.Client, error) {
	client, err := jsonrpc.Dial(RPC_TYPE, ip+":"+strconv.Itoa(port))
	if err != nil {
		log.Println(err)
	}
	return client, err
}
