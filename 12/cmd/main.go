package main

import (
	"flag"
	"fmt"
	"github.com/ilam072/wbtech-l2/12/internal/grep"
	"os"
)

func main() {
	var (
		after  = flag.Int("A", 0, "Print N lines after match")
		before = flag.Int("B", 0, "Print N lines before match")
		ctx    = flag.Int("C", 0, "Print N lines around match (both before and after)")

		countOnly  = flag.Bool("c", false, "Print only count of matching lines")
		ignoreCase = flag.Bool("i", false, "Ignore case distinctions")
		invert     = flag.Bool("v", false, "Invert match to select non-matching lines")
		fixed      = flag.Bool("F", false, "Interpret pattern as a fixed string (not a regex)")
		lineNums   = flag.Bool("n", false, "Prefix each line of output with the line number")
	)

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: grepx [OPTIONS] PATTERN [FILE]")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(2)
	}

	pattern := args[0]

	var in *os.File
	if len(args) >= 2 {
		f, err := os.Open(args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
		in = f
	} else {
		in = os.Stdin
	}

	opts := grep.Options{
		After:      *after,
		Before:     *before,
		Context:    *ctx, // передаём отдельное поле
		CountOnly:  *countOnly,
		IgnoreCase: *ignoreCase,
		Invert:     *invert,
		Fixed:      *fixed,
		LineNums:   *lineNums,
		Pattern:    pattern,
	}

	if err := grep.Run(in, os.Stdout, opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
