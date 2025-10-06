package shell

import (
	"fmt"
	"io"
	"strings"
)

func HandEcho(args []string, out io.Writer) {
	if len(args) > 0 {
		fmt.Fprintln(out, strings.Join(args, " "))
	} else {
		fmt.Fprintln(out)
	}
}
