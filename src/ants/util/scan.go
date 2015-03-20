package util

// scan spider from spiders dir
// for now it is useless,go cannot load libary dynamic

import (
	"go/parser"
	"go/token"
	"log"
)

func ScanSpider(dir string) {
	pkg, err := parser.ParseDir(token.NewFileSet(), dir, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(pkg)
	for k, v := range pkg {
		log.Println(k)
		for k1, v1 := range v.Files {
			log.Println(k1)
			for decls := range v1.Decls {
				log.Println(decls)
			}
		}
	}
}
