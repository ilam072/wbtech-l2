package main

import (
	"github.com/ilam072/wbtech-l2/13/internal/cut"
	"os"
)

func main() {
	cfg, err := cut.ParseFlags()
	if err != nil {
		cut.ExitWithError(err)
	}

	// Если указаны файлы
	if len(cfg.Files) > 0 {
		for _, filename := range cfg.Files {
			file, err := os.Open(filename)
			if err != nil {
				cut.ExitWithError(err)
			}
			err = cut.ProcessReader(file, os.Stdout, cfg)
			file.Close()
			if err != nil {
				cut.ExitWithError(err)
			}
		}
	} else {
		// STDIN
		err = cut.ProcessReader(os.Stdin, os.Stdout, cfg)
		if err != nil {
			cut.ExitWithError(err)
		}
	}
}
