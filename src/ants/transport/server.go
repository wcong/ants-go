package transport

import (
	"net"
)

type TcpServer struct {
	Linstener net.Listener
}

func NewTcpServer() *TcpServer {
	return &TcpServer{}
}

func (this *TcpServer) Listen() {

}
