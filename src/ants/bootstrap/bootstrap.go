package main

import (
	"ants/node"
	. "ants/util"
	"flag"
	"log"
	"os"
)

const (
	CONF_FILE = "/../conf/conf.json"
)

func initFlag(settings *Settings) {
	flag.IntVar(&settings.TcpPort, "tcp", settings.TcpPort, "tcp port")
	flag.IntVar(&settings.HttpPort, "http", settings.HttpPort, "http port")
}
func MakeSettings() *Settings {
	pwd, _ := os.Getwd()
	settings := LoadSettingFromFile(pwd + CONF_FILE)
	initFlag(settings)
	flag.Parse()
	return settings
}
func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("let us go shipping")
	setting := MakeSettings()
	Node := node.NewNode(setting)
	Node.Init()
	log.Println("finish init")
	Node.Start()
}
