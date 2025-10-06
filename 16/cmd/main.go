package main

import (
	"github.com/ilam072/wbtech-l2/16/internal/crawler"
	"github.com/ilam072/wbtech-l2/16/internal/flags"
)

func main() {
	f := flags.New()
	c := crawler.New(f.URL, f.Depth, f.Concurrency)
	c.Run()
}
