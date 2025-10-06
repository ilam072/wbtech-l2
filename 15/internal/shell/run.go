package shell

import (
	"fmt"
	"os"
	"os/exec"
)

func RunConditional(commands []ConditionalCommand) {
	lastStatus := 0
	for _, c := range commands {
		if c.Conditional == "&&" && lastStatus != 0 {
			break
		}
		if c.Conditional == "||" && lastStatus == 0 {
			break
		}

		err := runCommandWithRedirect(c.Cmd)

		if err != nil {
			fmt.Fprintf(os.Stderr, "minishell: %v\n", err)
			lastStatus = 1
		} else {
			lastStatus = 0
		}
	}
}

func runCommandWithRedirect(cmd Command) error {
	var stdin = os.Stdin
	var stdout = os.Stdout
	var err error

	var args []string
	i := 0
	for i < len(cmd.Args) {
		if cmd.Args[i] == ">" || cmd.Args[i] == ">>" || cmd.Args[i] == "<" {
			if i+1 >= len(cmd.Args) {
				return fmt.Errorf("syntax error near '%s'", cmd.Args[i])
			}
			filename := cmd.Args[i+1]
			switch cmd.Args[i] {
			case ">":
				stdout, err = os.Create(filename)
				if err != nil {
					return err
				}
			case ">>":
				stdout, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return err
				}
			case "<":
				stdin, err = os.Open(filename)
				if err != nil {
					return err
				}
			}
			i += 2
		} else {
			args = append(args, cmd.Args[i])
			i++
		}
	}

	var execErr error
	switch cmd.Name {
	case "pwd":
		execErr = HandPwd(stdout)
	case "echo":
		HandEcho(args, stdout)
	case "cd":
		execErr = HandCd(args)
	case "ps":
		execErr = HandPs(stdout)
	case "kill":
		execErr = HandKill(args)
	default:
		c := exec.Command(cmd.Name, args...)
		c.Stdin = stdin
		c.Stdout = stdout
		c.Stderr = os.Stderr
		execErr = c.Run()
		if err, ok := execErr.(*exec.ExitError); ok {
			_ = err
			execErr = nil
		}
	}

	if stdout != os.Stdout && stdout != nil {
		_ = stdout.Close()
	}
	if stdin != os.Stdin && stdin != nil {
		_ = stdin.Close()
	}

	return execErr
}
