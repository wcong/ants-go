package db

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestMysql(t *testing.T) {
	fmt.Println(sql.Drivers())
}
