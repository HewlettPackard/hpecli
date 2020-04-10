// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
	"github.com/HewlettPackard/hpecli/pkg/analytics"
	"github.com/HewlettPackard/hpecli/pkg/autocomplete"
	"github.com/HewlettPackard/hpecli/pkg/cloudvolume"
	"github.com/HewlettPackard/hpecli/pkg/greenlake"
	"github.com/HewlettPackard/hpecli/pkg/ilo"
	"github.com/HewlettPackard/hpecli/pkg/oneview"
	"github.com/HewlettPackard/hpecli/pkg/update"
	"github.com/HewlettPackard/hpecli/pkg/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// values are injected by the linker at build time via ldflags
var buildDate = "0"    // populated looks like: 2020-03-09
var gitCommitID = "0"  // populated looks like: fb75887
var sematicVer = "DEV" // populated looks like: 0.1.2

func main() {
	if err := run(); err != nil {
		fmt.Println("Error :", err)
		//nolint:gomnd // magic number ok here on exit
		os.Exit(1)
	}
}

func run() error {
	if isDebugLogging() {
		log.SetDebugLogging()
	}

	logrus.Debug("Started : Application initializing")
	defer logrus.Debug("Completed : Application shutdown")

	// create the root command.  It doesn't do anything, but used to hold
	// all of the other top level commands
	rootCmd := &cobra.Command{
		Use:           "hpe",
		Short:         "HPE Command Line Interface for accessing various services",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// this isn't acutually used.  It needs to be here so that when the
	// command line args are parsed, it won't complain that it is present
	// we already figured out if we are debugging via isDebugLogging
	_ = rootCmd.PersistentFlags().Bool("debug", false, "enable debug logging")

	// add all of the commands get added to this root one
	addSubCommands(rootCmd)

	// Are we are been called with no args at all?
	if len(os.Args) == 1 {
		// Then let's return the command usage as in: hpe --help
		os.Args = []string{os.Args[0], "--help"} 
	}  

	// execute the root command
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

func addSubCommands(rootCmd *cobra.Command) {
	vInfo := &version.Info{
		Sematic:     sematicVer,
		GitCommitID: gitCommitID,
		BuildDate:   buildDate,
		Verbose:     isDebugLogging(),
	}

	rootCmd.AddCommand(
		analytics.NewAnalyticsCommand(),
		autocomplete.NewAutoCompleteCommand(),
		cloudvolume.NewCloudVolumeCommand(),
		greenlake.NewGreenlakeCommand(),
		ilo.NewILOCommand(),
		oneview.NewOneViewCommand(),
		update.NewUpdateCommand(sematicVer),
		version.NewVersionCommand(vInfo),
	)
}

// cobra doesn't parse the command line arguments until cmd.Execute is called
// that is very late in the application initialization, so we have this simple
// check to see if they are requesting debug logging.  This allows us to set
// debug logging very early
func isDebugLogging() bool {
	for _, arg := range os.Args {
		if strings.EqualFold(arg, "--debug") {
			return true
		}
	}

	return false
}
