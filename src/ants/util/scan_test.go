package util

import (
	"os"
	"testing"
)

func TestScan(t *testing.T) {
	pwd, _ := os.Getwd()
	pwd += "/../../spiders"
	ScanSpider(pwd)
}
