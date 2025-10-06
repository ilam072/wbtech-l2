package main

import (
	"github.com/ilam072/wbtech-l2/10/cmd/root"
	"log"
)

func main() {
	if err := root.Execute(); err != nil {
		log.Fatalln(err)
	}
}
