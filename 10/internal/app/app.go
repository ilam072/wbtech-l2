package app

import (
	"github.com/ilam072/wbtech-l2/10/internal/flags"
	"github.com/ilam072/wbtech-l2/10/internal/printer"
	"github.com/ilam072/wbtech-l2/10/internal/sorter"
	"io"
	"log"
	"os"
)

type App struct {
	flags  flags.Flags
	sorter sorter.Sorter
	data   *os.File
}

func New(data *os.File, flags flags.Flags, sorter sorter.Sorter) *App {
	return &App{
		flags:  flags,
		sorter: sorter,
		data:   data,
	}
}

func (a *App) Run(output io.Writer) {
	a.sorter.ReadData(a.data)
	if a.flags.IgnoreTB {
		a.sorter.TrimTrailingSpaces()
	}

	if a.flags.Unique {
		a.sorter.Unique()
	}

	if a.flags.Column > 0 {
		a.sorter.SortByColumn(a.flags.Column, a.flags.Numeric)
	} else if a.flags.Numeric {
		a.sorter.SortNumeric()
	} else if a.flags.Month {
		a.sorter.SortByMonth(a.flags.Column)
	} else if a.flags.Human {
		a.sorter.SortHumanNumeric()
	} else if !a.flags.Check {
		a.sorter.Sort()
	}

	if a.flags.Reverse {
		a.sorter.Reverse()
	}

	if a.flags.Check {
		a.sorter.IsSorted(a.data.Name())
		return
	}

	if err := printer.Print(output, a.sorter.Data()); err != nil {
		log.Fatalf("failed to print data to output: %v", err.Error())
	}
}
