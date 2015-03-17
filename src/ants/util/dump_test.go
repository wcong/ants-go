package util

import (
	"testing"
)

func TestDump(t *testing.T) {
	DumpResult("../../../log", "test.log", "hello world")
}
