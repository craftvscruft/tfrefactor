package cmd

import (
	"fmt"
	"io/ioutil"
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
	flags.BoolP("force", "f", false, "Skip interactive approval of update before applying")

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
	if err != nil {
		return err
	}
	plan, err := refactor.Rename(fromAddress, toAddress, configPath)
	if err != nil {
		return err
	}
	err = approveAndApplyUpdate(cmd, plan)
	return err
}

func approveAndApplyUpdate(cmd *cobra.Command, plan *refactor.UpdatePlan) error {
	autoApprove, err := cmd.Flags().GetBool("force")
	CheckFatal(err)
	if len(plan.FileUpdates) > 0 {
		if autoApprove {
			err = applyUpdate(plan)
			CheckFatal(err)
			_, err = fmt.Fprintln(cmd.OutOrStdout(), "Done.")
			CheckFatal(err)
		} else {
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Update %v file(s)? [y/N]: ", len(plan.FileUpdates))
			CheckFatal(err)
			var in string
			_, _ = fmt.Fscanln(cmd.InOrStdin(), &in)
			// Ignore Fscanln err because empty input is OK.
			if in == "Y" || in == "y" {
				err = applyUpdate(plan)
				CheckFatal(err)
				_, err = fmt.Fprintln(cmd.OutOrStdout(), "Done.")
				CheckFatal(err)
			} else {
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "\nAborted.\n")
				CheckFatal(err)
			}
		}
	} else {
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "\nNo updates required.\n")
		CheckFatal(err)
	}
	return nil
}

func applyUpdate(plan *refactor.UpdatePlan) error {
	for _, update := range plan.FileUpdates {
		err := ioutil.WriteFile(update.Filename, []byte(update.AfterText), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
