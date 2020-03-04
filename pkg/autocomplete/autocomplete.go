// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package autocomplete

import (
	"github.com/HewlettPackard/hpecli/internal/platform/log"

	"github.com/spf13/cobra"
)

var genautocompleteCmd struct {
	autocompleteTarget string

	// bash for now (zsh and others will come)
	autocompleteType string
}

// Cmd represents the autocompletion command
var Cmd = &cobra.Command{
	Use:   "autocompletion",
	Short: "Generates bash completion scripts",
	Long: `Generates a shell autocompletion script for hpecli.

NOTE: The current version supports Bash only.
      This should work for *nix systems with Bash installed.

By default, the autocomplete.sh file will be generated in current folder

Add ` + "`--completionfile=/path/to/file`" + ` flag to set alternative
file-path and name.

Add the autocomplete.sh content to user ~/bash_profile file and reload the terminal`,
	RunE: run,
}

func run(cmd *cobra.Command, _ []string) error {
	if genautocompleteCmd.autocompleteType != "bash" {
		log.Logger.Warningf("Only Bash is supported for now")
		return nil
	}

	err := cmd.Root().GenBashCompletionFile(genautocompleteCmd.autocompleteTarget)

	if err != nil {
		return err
	}

	log.Logger.Debug("Bash completion file for hpecli saved to", genautocompleteCmd.autocompleteTarget)

	return nil
}

func init() {
	Cmd.PersistentFlags().StringVarP(&genautocompleteCmd.autocompleteTarget, "completionfile",
		"", "autocomplete.sh", "autocompletion file")
	Cmd.PersistentFlags().StringVarP(&genautocompleteCmd.autocompleteType, "type",
		"", "bash", "autocompletion type (currently only bash supported)")
}
