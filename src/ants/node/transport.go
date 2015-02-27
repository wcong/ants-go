package node

import (
	"ants/http"
	"ants/transport"
	"ants/util"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	TCP_EDN_SIGN         = "\t\n"
	TCP_EDN_SIGN_REPLACE = "\n"
	HADNLER_JOIN_REQUEST = iota
	HADNLER_JOIN_RESPONSE
	HANDLER_SEND_MASTER_REQUEST
)

// transport message struct
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
	Settings        *util.Settings
	TcpServer       net.Listener
	ConnMap         map[string]net.Conn
	HandleMap       map[int]func(*JSONMessage, net.Conn)
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
	handleMap := make(map[int]func(*JSONMessage, net.Conn))
	transporter := &Transporter{settings, ln, connMap, handleMap, "", node}
	transporter.HandleMap[HADNLER_JOIN_REQUEST] = transporter.handlerJoinRequest
	transporter.HandleMap[HANDLER_SEND_MASTER_REQUEST] = transporter.handlerSendMasterRequest
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
				jsonMessage := JSONMessage{
					Type:     HADNLER_JOIN_REQUEST,
					NodeInfo: *this.Node.NodeInfo,
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
			if redErr.Error() != "EOF" {
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
				log.Println(err)
			} else {
				go this.HandleMap[jsonMessage.Type](&jsonMessage, conn)
			}
		}
	}
}

// send message to node
func (this *Transporter) SendMessageToNode(nodeInfo *NodeInfo, message string) {
	this.SendMessage(this.ConnMap[nodeInfo.Name], message)
}

// send message by connection
// *		replace TCP_EDN_SIGN by TCP_EDN_SIGN_REPLACE
// *		send it
func (this *Transporter) SendMessage(conn net.Conn, message string) {
	if strings.Contains(message, TCP_EDN_SIGN) {
		message = strings.Replace(message, TCP_EDN_SIGN, TCP_EDN_SIGN_REPLACE, -1)
	}
	transport.SendMessage(conn, message+TCP_EDN_SIGN)
}

// what if some node what to join
func (this *Transporter) handlerJoinRequest(jsonMessage *JSONMessage, conn net.Conn) {
	if jsonMessage.NodeInfo.Ip != "" {
		log.Println("get node join request:ip:" + jsonMessage.NodeInfo.Ip + ";port:" + strconv.Itoa(jsonMessage.NodeInfo.Port))
		nodeInfo := &jsonMessage.NodeInfo
		if _, ok := this.ConnMap[nodeInfo.Name]; !ok {
			this.ConnMap[nodeInfo.Name] = conn
		}
		this.Node.AddNodeToCluster(nodeInfo)
	}
}

// deal with send master request,old master node elect new master node ,and send it to all node
func (this *Transporter) handlerSendMasterRequest(jsonMessage *JSONMessage, conn net.Conn) {
	this.Node.AddMasterNode(&jsonMessage.NodeInfo)
}
