package main

import (
	"fmt"
	eeaao_codegen "github.com/palindrom615/eeaao-codegen"
	"github.com/spf13/cobra"
	"os"
)

var (
	specdir    string
	codeletdir string
	outdir     string
	configFile string
)
var rootCmd = &cobra.Command{
	Use:   "eeaao-codegen-cli",
	Short: "anything code generator",
	Long: `anything code generator.
		You can generate anything from anything with this tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := eeaao_codegen.NewApp(specdir, outdir, codeletdir, configFile)
		app.Render()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&specdir, "specdir", "", "Directory for specifications")
	rootCmd.PersistentFlags().StringVar(&codeletdir, "codeletdir", "", "Directory for templates")
	rootCmd.PersistentFlags().StringVar(&outdir, "outdir", "build", "Directory for output")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file")

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
