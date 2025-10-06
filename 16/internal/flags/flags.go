package flags

import (
	"github.com/spf13/pflag"
)

type Flags struct {
	URL         string
	Depth       int
	Concurrency int
}

func New() *Flags {
	depth := pflag.IntP("depth", "d", 1, "глубина обхода сайта")
	concurrency := pflag.IntP("concurrency", "c", 5, "максимальное число одновременных загрузок")
	pflag.Parse()

	url := pflag.Arg(0)

	flags := Flags{
		URL:         url,
		Depth:       *depth,
		Concurrency: *concurrency,
	}

	return &flags
}
