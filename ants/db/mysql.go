package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type MysqlConnectMap struct {
	contectMap map[string]*sql.DB
}

func NewMysqlConnectMap() *MysqlConnectMap {
	conectMap := &MysqlConnectMap{}
	conectMap.InitMap()
	return conectMap
}

func (this *MysqlConnectMap) InitMap() {
	this.contectMap = make(map[string]*sql.DB)
}

func (this *MysqlConnectMap) InitConnection(spiderName, url string) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Println(err)
	}
	this.contectMap[spiderName] = db
}

func (this *MysqlConnectMap) CloseContection(spiderName string) {
	db, ok := this.contectMap[spiderName]
	if ok {
		db.Close()
		delete(this.contectMap, spiderName)
	}
}

func (this *MysqlConnectMap) Query(spiderName, query string, args ...interface{}) (*sql.Rows, error) {
	return this.contectMap[spiderName].Query(query, args...)
}

func (this *MysqlConnectMap) Exec(spiderName, query string, args ...interface{}) (sql.Result, error) {
	return this.contectMap[spiderName].Exec(query, args...)
}

var DefaultMysqlConnectMap = NewMysqlConnectMap()
