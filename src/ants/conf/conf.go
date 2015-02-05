package conf

import (
	"encoding/json"
	"os"
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
	return setting
}
