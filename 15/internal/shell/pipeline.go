package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

type Pipeline struct {
	Commands []Command
}

func ParsePipeline(line string) Pipeline {
	parts := strings.Split(line, "|")
	var cmds []Command
	for _, part := range parts {
		tokens := strings.Fields(strings.TrimSpace(part))
		if len(tokens) == 0 {
			continue
		}
		cmds = append(cmds, Command{
			Name: tokens[0],
			Args: tokens[1:],
		})
	}
	return Pipeline{Commands: cmds}
}

func handBuiltin(c Command, in io.Reader, out io.Writer) error {
	switch c.Name {
	case "pwd":
		return HandPwd(out)
	case "echo":
		if len(c.Args) > 0 {
			fmt.Fprintln(out, strings.Join(c.Args, " "))
		} else if in != nil {
			_, _ = io.Copy(out, in)
		} else {
			fmt.Fprintln(out)
		}
	case "cd":
		return HandCd(c.Args)
	case "ps":
		return HandPs(out)
	case "kill":
		return HandKill(c.Args)
	}
	return nil
}

func RunPipeline(p Pipeline) error {
	var prevReader *os.File
	var processes []*exec.Cmd
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	defer signal.Stop(sigs)

	for i, c := range p.Commands {
		isLast := i == len(p.Commands)-1
		var r, w *os.File
		var err error

		if !isLast {
			r, w, err = os.Pipe()
			if err != nil {
				return err
			}
		}

		if isBuiltin(c.Name) {
			out := os.Stdout
			if w != nil {
				out = w
			}
			if err := handBuiltin(c, prevReader, out); err != nil {
				if w != nil {
					w.Close()
				}
				if r != nil {
					r.Close()
				}
				return err
			}
			if w != nil {
				w.Close()
			}
			if prevReader != nil {
				prevReader.Close()
			}
			prevReader = r
		} else {
			cmd := exec.Command(c.Name, c.Args...)
			if prevReader != nil {
				cmd.Stdin = prevReader
			} else {
				cmd.Stdin = os.Stdin
			}
			if w != nil {
				cmd.Stdout = w
			} else {
				cmd.Stdout = os.Stdout
			}
			cmd.Stderr = os.Stderr

			if err := cmd.Start(); err != nil {
				if w != nil {
					w.Close()
				}
				if r != nil {
					r.Close()
				}
				return err
			}

			processes = append(processes, cmd)

			if prevReader != nil {
				prevReader.Close()
			}
			prevReader = r

			if w != nil {
				defer w.Close()
			}
		}
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for range sigs {
			for _, cmd := range processes {
				if cmd.Process != nil {
					_ = cmd.Process.Signal(syscall.SIGINT)
				}
			}
		}
	}()

	for _, cmd := range processes {
		_ = cmd.Wait()
	}

	close(sigs)
	<-done
	return nil
}

func IsPipeline(line string) bool {
	return strings.Contains(line, "|")
}

func isBuiltin(name string) bool {
	switch name {
	case "pwd", "echo", "cd", "ps", "kill":
		return true
	default:
		return false
	}
}

/*func handBuiltin(c Command, out io.Writer) error {
	var err error
	switch c.Name {
	case "pwd":
		err = HandPwd(out)
	case "echo":
		HandEcho(c.Args, out)
	case "cd":
		err = HandCd(c.Args)
	case "ps":
		err = HandPs(out)
	case "kill":
		err = HandKill(c.Args)
	}
	return err
}
*/
