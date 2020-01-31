// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/HewlettPackard/hpecli/pkg/store"
	"github.com/spf13/cobra"
)

var iloLoginData struct {
	host,
	username,
	password string
}

func init() {
	cmdILOLogin.Flags().StringVar(&iloLoginData.host, "host", "", "ilo ip address")
	cmdILOLogin.Flags().StringVarP(&iloLoginData.username, "username", "u", "", "ilo username")
	cmdILOLogin.Flags().StringVarP(&iloLoginData.password, "password", "p", "", "ilo passowrd")
}

// cmdIloLogin represents the get command
var cmdILOLogin = &cobra.Command{
	Use:   "login",
	Short: "Login to iLO: hpecli ilo login",
	RunE:  runILOLogin,
}

func runILOLogin(_ *cobra.Command, _ []string) error {
	if iloLoginData.host == "" {
		return fmt.Errorf("must provide --host or -h")
	}

	if !strings.HasPrefix(iloLoginData.host, "http") {
		iloLoginData.host = fmt.Sprintf("http://%s", iloLoginData.host)
	}

	if iloLoginData.username == "" {
		return fmt.Errorf("must provide --username or -u")
	}

	if iloLoginData.password == "" {
		// this really isn't secure to provide on the command line
		// need to provide a way to read from stdin
		return fmt.Errorf("must provide --password or -p")
	}

	logger.Debug("Attempting login with user: %v, at: %v", iloLoginData.username, iloLoginData.host)

	db, err := store.Open()
	if err != nil {
		logger.Debug("unable to open keystore: %v", err)
		return fmt.Errorf("%v", err)
	}
	defer db.Close()

	var val string
	if err := db.Get(key(), &val); err != nil {
		return fmt.Errorf("%v", err)
	}

	db.Close()

	// if err != nil {
	// 	logger.Debug("Unable to loginmust provide --password or -p because of: %v", err)
	// 	return fmt.Errorf("Unable to loginmust provide --password or -p")
	// }

	return nil
}

func key() string {
	return fmt.Sprintf("hpecli_ilo_token_%s", iloLoginData.host)
}
