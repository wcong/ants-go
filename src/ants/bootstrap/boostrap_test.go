package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestBootstrap(t *testing.T) {
	os.Args = []string{"tcp", "8200"}
	settings := MakeSettings()
	fmt.Println(settings.TcpPort)
	strings.NewReader("A")
}
