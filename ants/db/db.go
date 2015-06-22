package db

import (
	"database/sql"
)

/*
* contain all db connect for spiders
 */
type DbPool interface {
	InitMap()
	InitConnection(spiderName, url string)
	CloseContection(spiderName string)
	Query(spiderName, query string, args ...interface{}) (*sql.Rows, error)
	Exec(spiderName, query string, args ...interface{}) (sql.Result, error)
}
