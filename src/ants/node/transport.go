package node

import (
	"ants/transport"
	"ants/util"
	"encoding/json"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	TCP_EDN_SIGN         = "\t\n"
	TCP_EDN_SIGN_REPLACE = "\n"
)

// what a transporter do
// *		init conn
// *		send message
// *		accept message
// *		handle request
type Transporter struct {
	Settings        *util.Settings
	TcpServer       net.Listener
	ConnMap         map[string]net.Conn
	HandleMap       map[int]func(*RequestMessage, net.Conn)
	ServerTmpString string
	Node            *Node
}

// init a transporter
func NewTransporter(settings *util.Settings, node *Node) *Transporter {
	portString := strconv.Itoa(settings.TcpPort)
	ln, err := net.Listen("tcp", ":"+portString)
	log.Println("start to listen tcp:" + portString)
	if err != nil {
		panic(err)
	}
	connMap := make(map[string]net.Conn)
	handleMap := make(map[int]func(*RequestMessage, net.Conn))
	transporter := &Transporter{settings, ln, connMap, handleMap, "", node}
	transporter.HandleMap[HADNLER_JOIN_REQUEST] = transporter.handlerJoinRequest
	transporter.HandleMap[HADNLER_JOIN_RESPONSE] = transporter.handlerJoinResponse
	transporter.HandleMap[HANDLER_SEND_MASTER_REQUEST] = transporter.handlerSendMasterRequest
	transporter.HandleMap[HANDLER_SEND_REQUEST] = transporter.handlerSendRequest
	transporter.HandleMap[HANDLER_SEND_REQUEST_RESULT] = transporter.handlerSendRequestResult
	transporter.HandleMap[HANDLER_STOP_NODE] = transporter.handlerStopNode
	return transporter
}

// Transporter started
// loop tcp server
// connect to server and send join request
func (this *Transporter) Start() {
	go this.acceptRequest()
	if len(this.Settings.NodeList) > 0 {
		for _, nodeInfo := range this.Settings.NodeList {
			nodeSettings := strings.Split(nodeInfo, ":")
			ip := nodeSettings[0]
			port, _ := strconv.Atoi(nodeSettings[1])
			if ip == this.Node.NodeInfo.Ip && port == this.Node.NodeInfo.Port {
				continue
			}
			conn, err := transport.InitClient(ip, port)
			if err != nil {
				log.Println(err)
			} else {
				go this.ClientReader(conn)
				jsonMessage := RequestMessage{
					Type:     HADNLER_JOIN_REQUEST,
					NodeInfo: this.Node.NodeInfo,
				}
				message, _ := json.Marshal(jsonMessage)
				this.SendMessage(conn, string(message))
			}
		}
	}
}

// loop tcp server
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

// deap loop linsten client connection
func (this *Transporter) ClientReader(conn net.Conn) {
	buffer := make([]byte, 2048)
	for {
		_, redErr := conn.Read(buffer)
		if redErr != nil {
			time.Sleep(1 * time.Second)
			if redErr != io.EOF {
				log.Println(redErr)
			}
		} else {
			data := string(buffer)
			this.handleMessage(data, conn)
		}
	}
}

// when some message come ,we should deal with it
// *		cache  part message
// *		split by sign
// *		send it to handler by type of it
func (this *Transporter) handleMessage(data string, conn net.Conn) {
	log.Println("get data:" + data)
	if this.ServerTmpString != "" {
		data = this.ServerTmpString + data
		this.ServerTmpString = ""
	}
	if !strings.HasSuffix(data, TCP_EDN_SIGN) {
		lastIndex := strings.LastIndex(data, TCP_EDN_SIGN)
		if lastIndex >= 0 {
			this.ServerTmpString = data[lastIndex+len(TCP_EDN_SIGN):]
			data = data[:lastIndex]
		} else {
			this.ServerTmpString = data
			data = ""
		}
	} else {
		this.ServerTmpString = ""
	}
	if len(data) > 0 {
		splitString := strings.Split(data, TCP_EDN_SIGN)
		for _, jsonString := range splitString {
			var jsonMessage RequestMessage
			jsonString = strings.Trim(jsonString, "\x00")
			log.Println(jsonString)
			err := json.Unmarshal([]byte(jsonString), &jsonMessage)
			if err != nil {
				log.Panicln(err)
			} else {
				go this.HandleMap[jsonMessage.Type](&jsonMessage, conn)
			}
		}
	}
}

// send message to node
func (this *Transporter) SendMessageToNode(nodeName, message string) {
	this.SendMessage(this.ConnMap[nodeName], message)
}

// send message by connection
// *		replace TCP_EDN_SIGN by TCP_EDN_SIGN_REPLACE
// *		send it
func (this *Transporter) SendMessage(conn net.Conn, message string) {
	if strings.Contains(message, TCP_EDN_SIGN) {
		message = strings.Replace(message, TCP_EDN_SIGN, TCP_EDN_SIGN_REPLACE, -1)
	}
	log.Println("send message:" + message)
	transport.SendMessage(conn, message+TCP_EDN_SIGN)
}

// what if some node what to join
func (this *Transporter) handlerJoinRequest(jsonMessage *RequestMessage, conn net.Conn) {
	log.Println("get node join request:ip:" + jsonMessage.NodeInfo.Ip + ";port:" + strconv.Itoa(jsonMessage.NodeInfo.Port))
	nodeInfo := jsonMessage.NodeInfo
	this.ConnMap[nodeInfo.Name] = conn
	this.Node.AddNodeToCluster(nodeInfo)
	response := &RequestMessage{
		Type:     HADNLER_JOIN_RESPONSE,
		NodeInfo: this.Node.NodeInfo,
	}
	message, _ := json.Marshal(response)
	this.SendMessage(conn, string(message))
}

// deal with send master request,old master node elect new master node ,and send it to all node
func (this *Transporter) handlerSendMasterRequest(jsonMessage *RequestMessage, conn net.Conn) {
	this.Node.AddMasterNode(jsonMessage.NodeInfo)
	this.ConnMap[jsonMessage.NodeInfo.Name] = conn
}

func (this *Transporter) handlerJoinResponse(jsonMessage *RequestMessage, conn net.Conn) {
	this.ConnMap[jsonMessage.NodeInfo.Name] = conn
	this.Node.AddNodeToCluster(jsonMessage.NodeInfo)
}
func (this *Transporter) handlerSendRequest(jsonMessage *RequestMessage, conn net.Conn) {
	log.Println("get request" + jsonMessage.Request.UniqueName)
	this.Node.AcceptRequest(jsonMessage.Request)
}
func (this *Transporter) handlerSendRequestResult(jsonMessage *RequestMessage, conn net.Conn) {
	this.Node.AcceptResult(jsonMessage)
}
func (this *Transporter) handlerStopNode(jsonMessage *RequestMessage, conn net.Conn) {
	this.Node.StopCrawl()
}

// TODO when connect lost
