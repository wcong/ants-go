package transport

import (
	"log"
	"net"
	"strconv"
)

type TcpManager struct {
	TcpServer net.Listener
	ClientMap map[string]net.Conn
}

func handleRequest(conn net.Conn) {
	log.Println("get request")
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(buf))
}
func acceptRequest(ln net.Listener) {
	for {
		log.Println("loop")
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		handleRequest(conn)
	}
}

func InitServer(port int) {
	portString := strconv.Itoa(port)
	ln, err := net.Listen("tcp", ":"+portString)
	log.Println("start to listen tcp:" + portString)
	if err != nil {
		panic(err)
	}
	go acceptRequest(ln)
}
