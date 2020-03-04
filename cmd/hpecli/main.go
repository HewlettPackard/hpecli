// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
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

func main() {
	if err := run(); err != nil {
		fmt.Println("Error :", err)
		//nolint:gomnd // magic number ok here on exit
		os.Exit(1)
	}
}

func run() error {
	if isDebugLogging() {
		log.Logger.SetLevel(logrus.DebugLevel)
	}

	log.Logger.Debug("Started : Application initializing")
	defer log.Logger.Debug("Completed : Application shutdown")

	// channel to get async update
	isUpdateChan := make(chan bool)

	// async check if an update is available
	go func() {
		log.Logger.Debug("update : starting async check to see if an update is available")

		isUpdate := update.IsUpdateAvailable()

		log.Logger.Debugf("update : IsUpdateAvailable=%v", isUpdate)
		isUpdateChan <- isUpdate
	}()

	// create the root command.  It doesn't do anything, but used to hold
	// all of the other top level commands
	rootCmd := &cobra.Command{
		Use:           "hpecli",
		Short:         "hpe cli for accessing various services",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// this isn't acutually used.  It needs to be here so that when the
	// command line args are parsed, it won't complain that it is present
	// we already figured out if we are debugging via isDebugLogging
	_ = rootCmd.PersistentFlags().Bool("debug", false, "enable debug logging")

	// add all of the commands get added to this root one
	addSubCommands(rootCmd)

	// execute the root command
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	// check status from async method
	newRelease := <-isUpdateChan
	// update.Cmd.CalledAs() has a value if that was the command that was executed
	// if update was just run, we don't need to tell them that there is an update
	if newRelease && update.Cmd.CalledAs() == "" {
		log.Logger.Warn("An updated version of the CLI is available.  You can update by running \"hpecli update\"")
	}

	return nil
}

func addSubCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(cloudvolume.Cmd)
	rootCmd.AddCommand(ilo.Cmd)
	rootCmd.AddCommand(oneview.Cmd)
	rootCmd.AddCommand(update.Cmd)
	rootCmd.AddCommand(greenlake.Cmd)
	rootCmd.AddCommand(version.Cmd)
	rootCmd.AddCommand(autocomplete.Cmd)
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
