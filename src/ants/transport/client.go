package transport

import (
	"log"
	"net"
	"strconv"
)

func InitClient(ip string, port int) net.Conn {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
func SendMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
}
