package main

import (
	"bufio"
	"fmt"
	"github.com/ilam072/wbtech-l2/15/internal/shell"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	for {
		select {
		case <-sigs:
			fmt.Println()
			continue
		default:
			fmt.Print("> ")
			line, err := reader.ReadString('\n')
			if line == "" {
				continue
			}
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			if err == io.EOF {
				return
			}

			if shell.IsPipeline(line) {
				pipeline := shell.ParsePipeline(line)
				if err := shell.RunPipeline(pipeline); err != nil {
					fmt.Fprintf(os.Stderr, "minishell: %v\n", err)
				}
				continue
			}

			commands := shell.ParseConditional(line)
			if len(commands) == 0 {
				continue
			}

			shell.RunConditional(commands)
		}
	}
}
