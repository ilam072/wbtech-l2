package grep

import (
	"fmt"
	"io"
)

func Run(r io.Reader, out io.Writer, opts Options) error {
	// валидация
	if opts.After < 0 || opts.Before < 0 || opts.Context < 0 {
		return fmt.Errorf("context values must be non-negative")
	}

	m, err := newMatcher(opts.Pattern, opts.Fixed, opts.IgnoreCase)
	if err != nil {
		return fmt.Errorf("failed to build matcher: %w", err)
	}

	before := opts.Before
	if opts.Context > before {
		before = opts.Context
	}
	after := opts.After
	if opts.Context > after {
		after = opts.Context
	}

	lr := newLineReader(r)

	type entry struct {
		line string
		num  int
	}
	beforeBuf := make([]entry, 0, before)

	afterRemaining := 0
	matchCount := 0
	lastPrinted := 0

	for {
		line, ln, ok, err := lr.Next()
		if err != nil {
			return fmt.Errorf("failed to read line: %w", err)
		}
		if !ok {
			break
		}

		norm := normalizeLine(line, opts.Fixed, opts.IgnoreCase)
		match := m.Match(norm)
		if opts.Invert {
			match = !match
		}

		if opts.CountOnly {
			if match {
				matchCount++
			}
			if before > 0 {
				if len(beforeBuf) == before {
					// сдвигаем влево
					beforeBuf = beforeBuf[1:]
				}
				beforeBuf = append(beforeBuf, entry{line: line, num: ln})
			}
			continue
		}

		if match {
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
			beforeBuf = beforeBuf[:0]

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
			afterRemaining = after
			matchCount++
			continue
		}

		if afterRemaining > 0 {
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
		}

		if before > 0 {
			if len(beforeBuf) == before {
				beforeBuf = beforeBuf[1:]
			}
			beforeBuf = append(beforeBuf, entry{line: line, num: ln})
		}
	}

	if opts.CountOnly {
		_, err := fmt.Fprintln(out, matchCount)
		return err
	}

	return nil
}
