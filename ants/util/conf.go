package util

import (
	"encoding/json"
	"os"
	"strings"
)

// settings
type Settings struct {
	HttpPort         int
	MulticastEnable  bool
	Name             string
	NodeList         []string
	TcpPort          int
	LogPath          string
	ConfigFile       string
	DownloadInterval int
}

func NewSettings() *Settings {
	return &Settings{
		HttpPort:         8200,
		MulticastEnable:  false,
		Name:             "guess",
		NodeList:         []string{"127.0.0.1:8300"},
		TcpPort:          8300,
		LogPath:          "../log",
		DownloadInterval: 1,
	}
}

// load json config
// change 127.0.0.1 to basic ip
func LoadSettingFromFile(fileName string, setting *Settings) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(setting)
	if err != nil {
		panic(err)
	}
	for index, node := range setting.NodeList {
		nodeInfo := strings.Split(node, ":")
		if nodeInfo[0] != "127.0.0.1" {
			continue
		}
		newNodeInfo := GetLocalIp() + ":" + nodeInfo[1]
		setting.NodeList[index] = newNodeInfo
	}
}
