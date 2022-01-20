package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/craftvscruft/tfrefactor/refactor"
	"github.com/spf13/cobra"
)

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
		err := ioutil.WriteFile(update.Filename, []byte(update.AfterText), 0600)
		if err != nil {
			return err
		}
	}
	return nil
}
