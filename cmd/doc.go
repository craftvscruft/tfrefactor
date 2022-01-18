package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func init() {
	rootCmd.AddCommand(newDocCmd())
}

func newDocCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doc",
		Short: "Generate Markdown CLI docs in ./docs/cli",
		RunE:  runDocCmd,
	}
}

func runDocCmd(cmd *cobra.Command, args []string) error {
	return doc.GenMarkdownTree(rootCmd, "./docs/cli")
}
