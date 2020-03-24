// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package autocomplete

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type autocompleteOptions struct {
	targetFile string
	// bash for now (zsh and others will come)
	acType string
}

// NewAutoCompleteCommand to generate shell completion scripts for hpecli
func NewAutoCompleteCommand() *cobra.Command {
	var opts autocompleteOptions

	var cmd = &cobra.Command{
		Use:   "autocompletion",
		Short: "Generates a shell completion scripts for hpecli",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runAutoComplete(cmd, &opts)
		},
		Long: `Generates a shell completion scripts for hpecli.

	NOTE: The current version supports Bash, ZShell, and PowerShell.

	By default, the sh|zsh|ps file(s) will be generated in current folder

	Add ` + "`--completionfile=/path/to/file`" + ` flag to set alternative
	file-path and name.

	Add the generated autocomplete.sh content to user ~/bash_profile file and reload the terminal
	or just source the .sh file directly:`,
	}

	cmd.Flags().StringVarP(&opts.targetFile, "completionfile",
		"c", "", "autocompletion file")
	cmd.Flags().StringVarP(&opts.acType, "type",
		"t", "bash", "autocompletion shell type: {bash|zsh|powershell}")

	_ = cmd.MarkFlagRequired("type")

	return cmd
}
func runAutoComplete(cmd *cobra.Command, opts *autocompleteOptions) error {
	switch opts.acType {
	case "bash":
		if opts.targetFile == "" {
			opts.targetFile = "autocomplete.sh"
		}

		if err := cmd.Root().GenBashCompletionFile(opts.targetFile); err != nil {
			return err
		}
	case "zsh":
		if opts.targetFile == "" {
			opts.targetFile = "autocomplete.zsh"
		}

		if err := cmd.Root().GenZshCompletionFile(opts.targetFile); err != nil {
			return err
		}
	case "powershell":
		if opts.targetFile == "" {
			opts.targetFile = "autocomplete.ps1"
		}

		if err := cmd.Root().GenPowerShellCompletionFile(opts.targetFile); err != nil {
			return err
		}
	default:
		return logrus.Errorf("unsupported shell type %s", opts.acType)
	}

	logrus.Debugf("%s completion file for hpecli saved to %s", opts.acType, opts.targetFile)
	logrus.Printf(`%s completion file for hpecli saved in current folder as %s
	Add the generated file content to user ~/bash_profile file and reload the terminal or 
	just source the file directly`, opts.acType, opts.targetFile)

	return nil
}
