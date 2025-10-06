package grep

import (
	"fmt"
	"io"
)

// Run — потоковый grep с поддержкой контекста и флагов.
func Run(r io.Reader, out io.Writer, opts Options) error {
	// валидация
	if opts.After < 0 || opts.Before < 0 || opts.Context < 0 {
		return fmt.Errorf("context values must be non-negative")
	}

	// подготовить matcher
	m, err := newMatcher(opts.Pattern, opts.Fixed, opts.IgnoreCase)
	if err != nil {
		return fmt.Errorf("failed to build matcher: %w", err)
	}

	// эффективные before/after как у grep: берём максимум с Context
	before := opts.Before
	if opts.Context > before {
		before = opts.Context
	}
	after := opts.After
	if opts.Context > after {
		after = opts.Context
	}

	lr := newLineReader(r)

	// кольцевой буфер для before (хранит пары line, lineno)
	type entry struct {
		line string
		num  int
	}
	beforeBuf := make([]entry, 0, before)

	afterRemaining := 0
	matchCount := 0
	lastPrinted := 0 // номер последней распечатанной строки (чтобы не дублировать)

	for {
		line, ln, ok, err := lr.Next()
		if err != nil {
			return fmt.Errorf("failed to read line: %w", err)
		}
		if !ok {
			break
		}

		// нормализуем только для матчинга (в fixed + -i режиме)
		norm := normalizeLine(line, opts.Fixed, opts.IgnoreCase)
		match := m.Match(norm)
		if opts.Invert {
			match = !match
		}

		if opts.CountOnly {
			if match {
				matchCount++
			}
			// поддерживаем буфер перед (т.к. можем продолжать сканирование)
			if before > 0 {
				if len(beforeBuf) == before {
					// сдвигаем влево
					beforeBuf = beforeBuf[1:]
				}
				beforeBuf = append(beforeBuf, entry{line: line, num: ln})
			}
			// после не нужен в режиме -c
			continue
		}

		if match {
			// печатаем beforeBuf (только те, которые ещё не печатались)
			for i := 0; i < len(beforeBuf); i++ {
				e := beforeBuf[i]
				if e.num <= lastPrinted {
					continue
				}
				sep := "-"
				if opts.LineNums {
					if _, err := fmt.Fprintf(out, "%d%s%s\n", e.num, sep, e.line); err != nil {
						return err
					}
				} else {
					if _, err := fmt.Fprintf(out, "%s\n", e.line); err != nil {
						return err
					}
				}
				lastPrinted = e.num
			}
			// очистить буфер before (мы уже вывели его)
			beforeBuf = beforeBuf[:0]

			// печатаем совпадающую строку с ':'
			sep := ":"
			if opts.LineNums {
				if _, err := fmt.Fprintf(out, "%d%s%s\n", ln, sep, line); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprintf(out, "%s\n", line); err != nil {
					return err
				}
			}
			lastPrinted = ln
			// включаем печать after
			afterRemaining = after
			matchCount++
			continue
		}

		// если сейчас находимся в режиме печати after (последующие строки после совпадения)
		if afterRemaining > 0 {
			// печатаем как контекстную строку с '-'
			sep := "-"
			if opts.LineNums {
				if _, err := fmt.Fprintf(out, "%d%s%s\n", ln, sep, line); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprintf(out, "%s\n", line); err != nil {
					return err
				}
			}
			lastPrinted = ln
			afterRemaining--
			// и при печати after мы всё равно добавляем текущую строку в beforeBuf ниже
		}

		// сохраняем текущую строку в beforeBuf (для будущих совпадений)
		if before > 0 {
			if len(beforeBuf) == before {
				beforeBuf = beforeBuf[1:]
			}
			beforeBuf = append(beforeBuf, entry{line: line, num: ln})
		}
	}

	// в конце, если -c, вернуть количество
	if opts.CountOnly {
		_, err := fmt.Fprintln(out, matchCount)
		return err
	}

	return nil
}
