package util

import (
	"encoding/json"
	"os"
	"strings"
)

type Settings struct {
	HttpPort        int
	MulticastEnable bool
	Name            string
	NodeList        []string
	TcpPort         int
}

func LoadSettingFromFile(fileName string) *Settings {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	setting := &Settings{}
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
	return setting
}
