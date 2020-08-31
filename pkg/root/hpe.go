package root

import (
	"github.com/HewlettPackard/hpecli/pkg/analytics"
	"github.com/HewlettPackard/hpecli/pkg/autocomplete"
	"github.com/HewlettPackard/hpecli/pkg/cloudvolume"
	"github.com/HewlettPackard/hpecli/pkg/ilo"
	"github.com/HewlettPackard/hpecli/pkg/oneview"
	"github.com/HewlettPackard/hpecli/pkg/update"
	"github.com/HewlettPackard/hpecli/pkg/version"
	"github.com/spf13/cobra"
)




// create the root command.  It doesn't do anything, but used to hold
// all of the other top level commands
func NewRootCommand(version *version.Info) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "hpe",
		Short:         "HPE Command Line Interface for accessing various services",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	addSubCommands(cmd, version)

	// this isn't acutually used.  It needs to be here so that when the
	// command line args are parsed, it won't complain that it is present
	// we already figured out if we are debugging via isDebugLogging
	_ = cmd.PersistentFlags().Bool("debug", false, "enable debug logging")
	//_ = cmd.PersistentFlags().Bool("debugHttp", false, "enable http debug logging - may print confidential information")

	return cmd
}

func addSubCommands(rootCmd *cobra.Command, vInfo *version.Info) {
	rootCmd.AddCommand(
		analytics.NewAnalyticsCommand(),
		autocomplete.NewAutoCompleteCommand(),
		cloudvolume.NewCloudVolumeCommand(),
		ilo.NewILOCommand(),
		oneview.NewOneViewCommand(),
		update.NewUpdateCommand(vInfo.Sematic),
		version.NewVersionCommand(vInfo),
	)
}