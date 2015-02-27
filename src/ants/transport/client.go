package transport

import (
	"log"
	"net"
	"strconv"
)

func InitClient(ip string, port int) (net.Conn, error) {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	return conn, err
}
func SendMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Println(err)
	}
}
