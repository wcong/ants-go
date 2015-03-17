package util

import (
	"log"
	"os"
	"time"
)

// os.Getwn() + "/"+path
// add fileName with message
func DumpResult(path, spiderName, message string) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	makePath := pwd + "/" + path
	MkdirIfNotExist(makePath)
	fileName := time.Now().Format("2006-01-02T15:04:05") + "-" + spiderName + ".log"
	file, fileErr := os.Create(makePath + "/" + fileName)
	if fileErr != nil {
		log.Println(fileErr)
		return
	}
	defer file.Close()
	file.Write([]byte(message))
}

func MkdirIfNotExist(path string) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			makeDirErr := os.Mkdir(path, os.ModePerm)
			if makeDirErr != nil {
				log.Println(makeDirErr)
				return
			}
		}
	}
}
