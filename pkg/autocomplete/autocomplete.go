// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package autocomplete

import (
	"errors"

	"github.com/HewlettPackard/hpecli/internal/platform/log"

	"github.com/spf13/cobra"
)

type autocompleteOptions struct {
	targetFile string
	// bash for now (zsh and others will come)
	acType string
}

func NewAutoCompleteCommand() *cobra.Command {
	var opts autocompleteOptions

	var cmd = &cobra.Command{
		Use:   "autocompletion",
		Short: "Generates bash completion scripts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runAutoComplete(cmd, &opts)
		},
		Long: `Generates a shell autocompletion script for hpecli.

	NOTE: The current version supports Bash only.
		  This should work for *nix systems with Bash installed.

	By default, the autocomplete.sh file will be generated in current folder

	Add ` + "`--completionfile=/path/to/file`" + ` flag to set alternative
	file-path and name.

	Add the autocomplete.sh content to user ~/bash_profile file and reload the terminal`,
	}

	cmd.PersistentFlags().StringVarP(&opts.targetFile, "completionfile",
		"", "autocomplete.sh", "autocompletion file")
	cmd.PersistentFlags().StringVarP(&opts.acType, "type",
		"", "bash", "autocompletion type (currently only bash supported)")

	return cmd
}

func runAutoComplete(cmd *cobra.Command, opts *autocompleteOptions) error {
	if opts.acType != "bash" {
		return errors.New("bash is the only currently supported shell type")
	}

	if err := cmd.Root().GenBashCompletionFile(opts.targetFile); err != nil {
		return err
	}

	log.Logger.Debug("Bash completion file for hpecli saved to", opts.targetFile)

	return nil
}
