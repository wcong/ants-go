package node

import (
	"ants/conf"
	"ants/http"
	"ants/transport"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	TCP_EDN_SIGN         = "\t\n"
	TCP_EDN_SIGN_REPLACE = "\n"
	HADNLE_JOIN_REQUEST  = iota
)

type JSONMessage struct {
	Type     int
	Request  http.Request
	NodeInfo NodeInfo
}

// what a transporter do
// *		init conn
// *		send message
// *		accept message
// *		handle request
type Transporter struct {
	Settings        *conf.Settings
	TcpServer       net.Listener
	ConnMap         map[*NodeInfo]net.Conn
	HandleMap       map[int]func(*JSONMessage, net.Conn)
	ServerTmpString string
	Node            *Node
}

func NewTransporter(settings *conf.Settings, node *Node) *Transporter {
	portString := strconv.Itoa(settings.TcpPort)
	ln, err := net.Listen("tcp", ":"+portString)
	log.Println("start to listen tcp:" + portString)
	if err != nil {
		panic(err)
	}
	connMap := make(map[*NodeInfo]net.Conn)
	handleMap := make(map[int]func(*JSONMessage, net.Conn))
	transporter := &Transporter{settings, ln, connMap, handleMap, "", node}
	transporter.HandleMap[HADNLE_JOIN_REQUEST] = transporter.handleJoinRequest
	return transporter
}

func (this *Transporter) acceptRequest() {
	for {
		log.Println("loop")
		conn, err := this.TcpServer.Accept()
		if err != nil {
			log.Fatal(err)
		}
		buf := make([]byte, 1024)
		_, redErr := conn.Read(buf)
		if redErr != nil {
			log.Fatal(redErr)
		}
		data := string(buf)
		this.handleMessage(data, conn)
	}
}
func (this *Transporter) Start() {
	go this.acceptRequest()
	if len(this.Settings.NodeList) > 0 {
		for _, nodeInfo := range this.Settings.NodeList {
			nodeSettings := strings.Split(nodeInfo, ":")
			port, _ := strconv.Atoi(nodeSettings[1])
			client := transport.InitClient(nodeSettings[0], port)
			go this.ClientReader(client)
		}
	}
}
func (this *Transporter) handleJoinRequest(jsonMessage *JSONMessage, conn net.Conn) {
	if jsonMessage.NodeInfo.Ip != "" {
		nodeInfo := &jsonMessage.NodeInfo
		if _, ok := this.ConnMap[nodeInfo]; !ok {
			this.ConnMap[nodeInfo] = conn
		}
		this.Node.AddNodeToCluster(nodeInfo)
	}
}
func (this *Transporter) ClientReader(conn net.Conn) {
	buffer := make([]byte, 2048)
	for {
		_, redErr := conn.Read(buffer)
		if redErr != nil {
			log.Fatal(redErr)
		}
		data := string(buffer)
		this.handleMessage(data, conn)
	}
}
func (this *Transporter) handleMessage(data string, conn net.Conn) {
	if this.ServerTmpString != "" {
		data += this.ServerTmpString
		this.ServerTmpString = ""
	}
	if !strings.HasSuffix(data, TCP_EDN_SIGN) {
		lastIndex := strings.LastIndex(data, TCP_EDN_SIGN)
		if lastIndex >= 0 {
			this.ServerTmpString = data[lastIndex:]
			data = data[:lastIndex]
		} else {
			this.ServerTmpString = data
			data = ""
		}
	}
	if len(data) > 0 {
		splitString := strings.Split(data, TCP_EDN_SIGN)
		for _, jsonString := range splitString {
			var jsonMessage JSONMessage
			err := json.Unmarshal([]byte(jsonString), &jsonMessage)
			if err != nil {
				log.Fatal(err)
			} else {
				go this.HandleMap[jsonMessage.Type](&jsonMessage, conn)
			}
		}
	}
}

func (this *Transporter) SendMessage(nodeInfo *NodeInfo, message string) {
	if strings.Contains(message, TCP_EDN_SIGN) {
		message = strings.Replace(message, TCP_EDN_SIGN, TCP_EDN_SIGN_REPLACE, -1)
	}
	transport.SendMessage(this.ConnMap[nodeInfo], message)
}
