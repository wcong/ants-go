package main

import (
	. "ants/conf"
	"flag"
)

func initFlag() *Settings {
	settings := &Settings{1}
	flag.IntVar(&settings.TcpPort, "tcp", 8200, "tcp port")
	return settings
}
func MakeSettings() *Settings {
	settings := initFlag()
	flag.Parse()
	return settings
}
