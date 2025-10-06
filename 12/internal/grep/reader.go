package grep

import (
	"bufio"
	"io"
)

type lineReader struct {
	scanner *bufio.Scanner
	lineNum int
}

func newLineReader(r io.Reader) *lineReader {
	sc := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 1<<20)
	return &lineReader{scanner: sc}
}

func (lr *lineReader) Next() (string, int, bool, error) {
	if lr.scanner.Scan() {
		lr.lineNum++
		return lr.scanner.Text(), lr.lineNum, true, nil
	}
	if err := lr.scanner.Err(); err != nil {
		return "", lr.lineNum, false, err
	}
	return "", lr.lineNum, false, nil
}
