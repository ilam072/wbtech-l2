package cut

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Fields    []int
	Delimiter string
	Separated bool
	Files     []string // новые аргументы файлов
}

func ParseFlags() (*Config, error) {
	fieldsStr := flag.String("f", "", "fields to output (e.g. 1,3-5)")
	delimiter := flag.String("d", "\t", "field delimiter")
	separated := flag.Bool("s", false, "only lines with delimiter")
	flag.Parse()

	if *fieldsStr == "" {
		return nil, fmt.Errorf("error: -f flag is required")
	}

	fields, err := ParseFields(*fieldsStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing fields: %w", err)
	}

	files := flag.Args()

	return &Config{
		Fields:    fields,
		Delimiter: *delimiter,
		Separated: *separated,
		Files:     files,
	}, nil
}

func ExitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
