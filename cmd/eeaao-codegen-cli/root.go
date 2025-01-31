package main

import (
	"fmt"
	eeaao_codegen "github.com/palindrom615/eeaao-codegen"
	"github.com/spf13/cobra"
	"os"
)

var (
	codeletdir string
	outdir     string
	valuesFile string
)
var rootCmd = &cobra.Command{
	Use:   "eeaao-codegen-cli",
	Short: "anything code generator",
	Long: `anything code generator.
		You can generate anything from anything with this tool.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app := eeaao_codegen.NewApp(outdir, codeletdir, valuesFile)
		return app.Render()
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&codeletdir, "codeletdir", "", "Directory for templates")
	rootCmd.PersistentFlags().StringVar(&outdir, "outdir", "build", "Directory for output")
	rootCmd.PersistentFlags().StringVar(&valuesFile, "value", "", "value file")

	rootCmd.MarkPersistentFlagRequired("codeletdir")
	rootCmd.MarkPersistentFlagRequired("outdir")

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
