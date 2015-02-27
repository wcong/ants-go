package util

import (
	"fmt"
	"os"
	"testing"
)

func TestConf(t *testing.T) {
	pwd, _ := os.Getwd()
	pwd += "/../../../conf/conf.json"
	setting := LoadSettingFromFile(pwd)
	fmt.Println(setting.HttpPort)
}
