package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is a top level command instance
var rootCmd = &cobra.Command{
	Use:           "tfrefactor",
	Short:         "Automated refactoring for Terraform",
	SilenceErrors: false,
	SilenceUsage:  true,
}

func init() {
	setDefaultStream(rootCmd)
}

func setDefaultStream(cmd *cobra.Command) {
	cmd.SetIn(os.Stdin)
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
}

func Execute() {
	rootCmd.Execute()
}
