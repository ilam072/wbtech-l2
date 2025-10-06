package cut

import (
	"bufio"
	"io"
	"strings"
)

// ProcessReader читает строки из r, обрабатывает их и выводит в w
func ProcessReader(r io.Reader, w io.Writer, cfg *Config) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		output := line // по умолчанию выводим всю строку

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
				// Если флаг -s и поле отсутствует → пропускаем строку
				continue
			}
		} else if cfg.Separated {
			// Если строка не содержит разделителя и -s указан → пропускаем
			continue
		}

		// Всегда добавляем LF после строки (совпадение с cut)
		w.Write([]byte(output + "\n"))
	}

	return scanner.Err()
}
