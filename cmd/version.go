package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "0.0.1"
)

func init() {
	rootCmd.AddCommand(newVersionCmd())
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		RunE:  runVersionCmd,
	}
}

func runVersionCmd(cmd *cobra.Command, args []string) error {
	_, err := fmt.Fprintln(cmd.OutOrStdout(), Version)
	return err
}
