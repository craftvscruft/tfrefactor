package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// RootCmd is a top level command instance
var rootCmd = &cobra.Command{
	Use:           "tf-rf",
	Short:         "Automated refactoring for Terraform",
	SilenceErrors: true,
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
