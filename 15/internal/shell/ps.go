package shell

import (
	"fmt"
	"github.com/shirou/gopsutil/v4/process"
	"io"
	"strings"
)

func HandPs(out io.Writer) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%8s %8s %10s\n", "PID", "PPID", "NAME"))

	processes, err := process.Processes()
	if err != nil {
		return err
	}

	for _, p := range processes {
		pid := p.Pid
		ppid, _ := p.Ppid()
		name, _ := p.Name()
		sb.WriteString(fmt.Sprintf("%8d %8d %10s\n", pid, ppid, name))
	}

	fmt.Fprint(out, sb.String())
	return nil
}
