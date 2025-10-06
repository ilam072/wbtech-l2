package shell

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func HandKill(args []string) error {
	if len(args) < 1 {
		return errors.New("usage: kill pid")
	}

	for _, a := range args {
		pid, err := strconv.Atoi(a)
		if err != nil {
			return fmt.Errorf("%s: arguments must be pid", a)
		}

		proc, err := os.FindProcess(pid)
		if err != nil {
			return fmt.Errorf("(%v) - No such process", pid)
		}

		err = proc.Kill()
		if err != nil {
			if errors.Is(err, os.ErrPermission) {
				return fmt.Errorf("(%v) - Operation not permitted", pid)
			}
			if errors.Is(err, os.ErrProcessDone) {
				return fmt.Errorf("(%v) - No such process", pid)
			}
			return fmt.Errorf("(%v) - %v", pid, err)
		}
	}

	return nil
}
