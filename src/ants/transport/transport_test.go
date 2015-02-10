package transport

import (
	"log"
	"testing"
	"time"
)

func TestTransport(t *testing.T) {
	var port int = 8300
	InitServer(port)
	time.Sleep(3 * time.Second)
	log.Println("client")
	conn := InitClient("127.0.0.1", port)
	SendMessage(conn, "hello world")
	conn.Close()
	time.Sleep(3 * time.Second)
	log.Println("close")
}
