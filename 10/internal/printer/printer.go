package printer

import (
	"fmt"
	"io"
)

func Print(output io.Writer, data []string) error {
	for _, line := range data {
		_, err := fmt.Fprint(output, line)
		if err != nil {
			return err
		}
	}
	return nil
}
