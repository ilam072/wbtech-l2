package shell

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// todo: использовать функцию
func External(command Command) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	cmd := exec.Command(command.Name, command.Args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	go func() {
		for range sigs {
			if cmd.Process != nil {
				_ = cmd.Process.Signal(syscall.SIGINT)
			}
		}
	}()

	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil
		}
		return fmt.Errorf("%s: command not found", command.Name)
	}

	return nil
}
