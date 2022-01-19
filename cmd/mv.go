package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/craftvscruft/tfrefactor/refactor"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newMvCmd())
}

func newMvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mv <ADDRESS> <TO_FILE>",
		Short: "Move element to a different file",
		Long: `Move element to a different file.
Can be var, data, resource, output, local, or module.

Arguments:
  ADDRESS     The address (e.g. var.a, data.vpc.default, aws_vpc.default).
  TO_FILE     File to move to.
`,
		RunE: runMvCmd,
	}
	flags := cmd.Flags()
	flags.StringP("config", "c", "-", "Path of terraform to modify, defaults to current.")
	flags.BoolP("force", "f", false, "Skip interactive approval of update before applying")

	return cmd
}

func runMvCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 argument, but got %d arguments", len(args))
	}

	fromAddr := args[0]
	toFile := args[1]

	configPath, err := cmd.Flags().GetString("config")
	CheckFatal(err)

	if configPath == "-" {
		configPath, err = os.Getwd()
		CheckFatal(err)
	}

	if !path.IsAbs(toFile) {
		toFile = path.Join(configPath, toFile)
	}
	plan, err := refactor.Mv(fromAddr, toFile, configPath)
	if err != nil {
		return err
	}
	err = approveAndApplyUpdate(cmd, plan)
	return err
}
