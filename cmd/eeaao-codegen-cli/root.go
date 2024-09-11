package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	specdir    string
	codeletdir string
	outdir     string
)
var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("specdir: %s\n", specdir)
		fmt.Printf("codeletdir: %s\n", codeletdir)
		fmt.Printf("outdir: %s\n", outdir)
		// Do Stuff Here
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&specdir, "specdir", "", "Directory for specifications")
	rootCmd.PersistentFlags().StringVar(&codeletdir, "codeletdir", "", "Directory for templates")
	rootCmd.PersistentFlags().StringVar(&outdir, "outdir", "", "Directory for output")

	rootCmd.MarkPersistentFlagRequired("specdir")
	rootCmd.MarkPersistentFlagRequired("codeletdir")
	rootCmd.MarkPersistentFlagRequired("outdir")

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
