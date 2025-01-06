package main

import (
	eeaao_codegen "github.com/palindrom615/eeaao-codegen"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(shellCmd)
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Start a REPL shell for testing render.star",
	Run: func(cmd *cobra.Command, args []string) {
		app := eeaao_codegen.NewApp(outdir, codeletdir, valuesFile)
		app.RunShell()
	},
}
