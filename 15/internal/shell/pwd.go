package shell

import (
	"fmt"
	"io"
	"os"
)

func HandPwd(out io.Writer) error {
	curDir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Fprintln(out, curDir)
	return nil
}
