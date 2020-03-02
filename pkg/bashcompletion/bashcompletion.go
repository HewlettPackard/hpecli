// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package bashcompletion

import (
	"hpecli/pkg/logger"

	"github.com/spf13/cobra"
)

type genautocompleteCmd struct {
	autocompleteTarget string

	// bash for now (zsh and others will come)
	autocompleteType string
}

// Cmd represents the completion command
var Cmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `Generates a shell autocompletion script for hpecli.

NOTE: The current version supports Bash only.
      This should work for *nix systems with Bash installed.

By default, the file is written directly to /etc/bash_completion.d
for convenience, and the command may need superuser rights, e.g.:

	$ sudo hpecli completion

Add ` + "`--completionfile=/path/to/file`" + ` flag to set alternative
file-path and name.

Logout and in again to reload the completion scripts,
or just source them in directly:

	$ . /etc/bash_completion`,
	Run: run,
}

func run(cmd *cobra.Command, _ []string) {
	// Cmd.GenBashCompletion(os.Stdout)
	cmd.Root().GenBashCompletionFile("pramod.sh")
	logger.Debug("Bash completion file for hpecli saved to", genautocompleteCmd.autocompleteTarget)
}

func init() {
	Cmd.PersistentFlags().StringVarP(&genautocompleteCmd.autocompleteTarget, "completionfile", "", "/etc/bash_completion.d/hpecli.sh", "autocompletion file")
	Cmd.PersistentFlags().StringVarP(&genautocompleteCmd.autocompleteType, "type", "", "bash", "autocompletion type (currently only bash supported)")

	// For bash-completion
	Cmd.PersistentFlags().SetAnnotation("completionfile", cobra.BashCompFilenameExt, []string{})
}
