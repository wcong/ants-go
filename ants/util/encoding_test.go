package util

import (
	"fmt"
	"strconv"
	"testing"
)

func TestEncoding(t *testing.T) {
	a := HashString("127.0.0.1" + strconv.Itoa(9300))
	fmt.Println(a)
	b := strconv.FormatUint(a, 10)
	fmt.Println(string(b))
}
