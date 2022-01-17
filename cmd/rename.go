package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/craftvscruft/tfrefactor/refactor"
)

func init() {
	rootCmd.AddCommand(newRenameCmd())
}

func newRenameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename <FROM_ADDRESS> <TO_ADDRESS>",
		Short: "Rename var / data / resource",
		Long: `Rename var / data / resource

Arguments:
  FROM_ADDRESS     The old address (e.g. var.a, data.vpc.default, aws_vpc.default).
  TO_ADDRESS       The new address.
`,
		RunE: runRenameCmd,
	}
	flags := cmd.Flags()
	flags.StringP("config", "c", "-", "Path of terraform to modify, defaults to current.")

	return cmd
}

func CheckFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func runRenameCmd(cmd *cobra.Command, args []string) error {
	setDefaultStream(cmd)
	if len(args) != 2 {
		return fmt.Errorf("expected 2 argument, but got %d arguments", len(args))
	}

	fromAddress := args[0]
	toAddress := args[1]
	configPath, err := cmd.Flags().GetString("config")
	CheckFatal(err)
	if configPath == "-" {
		configPath, err = os.Getwd()
		CheckFatal(err)
	}

	CheckFatal(err)
	_, err = fmt.Fprintf(cmd.OutOrStdout(), "Renaming '%v' -> '%v' in %v\n", fromAddress, toAddress, configPath)
	return refactor.Rename(fromAddress, toAddress, configPath)
}
