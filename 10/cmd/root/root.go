package root

import (
	"github.com/ilam072/wbtech-l2/10/internal/app"
	"github.com/ilam072/wbtech-l2/10/internal/flags"
	"github.com/ilam072/wbtech-l2/10/internal/sorter"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var f = flags.New()

var rootCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sort lines",
	Run: func(cmd *cobra.Command, args []string) {
		data := os.Stdin
		if len(args) > 0 {
			file, err := os.Open(args[0])
			if err != nil {
				log.Fatalln("failed to open file: %w", err)
			}
			data = file

			defer func() {
				if err := data.Close(); err != nil {
					log.Fatalln(err)
				}
			}()
		}

		s := sorter.New()
		a := app.New(data, f, s)
		a.Run(os.Stdout)
	},
}

func init() {
	rootCmd.Flags().IntVarP(&f.Column, "key", "k", 0, "Sort by column N")
	rootCmd.Flags().BoolVarP(&f.Numeric, "numeric", "n", false, "Numeric sort")
	rootCmd.Flags().BoolVarP(&f.Reverse, "reverse", "r", false, "Reverse order")
	rootCmd.Flags().BoolVarP(&f.Unique, "unique", "u", false, "Unique lines")
	rootCmd.Flags().BoolVarP(&f.Month, "month", "M", false, "Sort by month name")
	rootCmd.Flags().BoolVarP(&f.IgnoreTB, "ignore-trailing-blanks", "b", false, "Ignore trailing blanks")
	rootCmd.Flags().BoolVarP(&f.Check, "check", "c", false, "Check if sorted")
	rootCmd.Flags().BoolVarP(&f.Human, "human-numeric-sort", "H", false, "Human-readable numeric sort")
}

func Execute() error {
	return rootCmd.Execute()
}
