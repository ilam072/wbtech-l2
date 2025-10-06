package cut

import (
	"bufio"
	"io"
	"strings"
)

func ProcessReader(r io.Reader, w io.Writer, cfg *Config) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		output := line

		if strings.Contains(line, cfg.Delimiter) {
			cols := strings.Split(line, cfg.Delimiter)
			outCols := make([]string, 0, len(cfg.Fields))

			for _, idx := range cfg.Fields {
				if idx >= 0 && idx < len(cols) {
					outCols = append(outCols, cols[idx])
				}
			}

			if len(outCols) > 0 {
				output = strings.Join(outCols, cfg.Delimiter)
			} else if cfg.Separated {
				continue
			}
		} else if cfg.Separated {
			continue
		}

		w.Write([]byte(output + "\n"))
	}

	return scanner.Err()
}
