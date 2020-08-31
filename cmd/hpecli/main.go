// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
	"github.com/HewlettPackard/hpecli/pkg/root"
	"github.com/HewlettPackard/hpecli/pkg/version"
	"github.com/sirupsen/logrus"
)

// values are injected by the linker at build time via ldflags
var buildDate = "0"    // populated looks like: 2020-03-09
var gitCommitID = "0"  // populated looks like: fb75887
var semanticVer = "DEV" // populated looks like: 0.1.2

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
//	if isDebugHttpLogging() {
//		utils.DebugLogHttpCalls()
//	}

	logrus.Debug("Started : Application initializing")
	defer logrus.Debug("Completed : Application shutdown")

	vInfo := &version.Info{
		Sematic:     semanticVer,
		GitCommitID: gitCommitID,
		BuildDate:   buildDate,
		Verbose:     isDebugLogging(),
	}

	root := root.NewRootCommand(vInfo)


	// Are we are been called with no args at all?
	if len(os.Args) == 1 {
		// Then let's return the command usage as in: hpe --help
		os.Args = []string{os.Args[0], "--help"} 
	}  

	// execute the root command
	if err := root.Execute(); err != nil {
		return err
	}

	return nil
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

// cobra doesn't parse the command line arguments until cmd.Execute is called
// that is very late in the application initialization, so we have this simple
// check to see if they are requesting debug logging.  This allows us to set
// debug logging very early
func isDebugHTTPLogging() bool {
	for _, arg := range os.Args {
		if strings.EqualFold(arg, "--debugHttp") {
			return true
		}
	}

	return false
}