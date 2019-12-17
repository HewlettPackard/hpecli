// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package ilo

import (
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/HewlettPackard/hpecli/pkg/store"
	"github.com/spf13/cobra"
)

var (
	host     string
	username string
	password string
)

func init() {
	cmdIloLogin.Flags().StringVar(&host, "host", "", "ilo ip address")
	cmdIloLogin.Flags().StringVarP(&username, "username", "u", "", "ilo username")
	cmdIloLogin.Flags().StringVarP(&password, "password", "p", "", "ilo passowrd")
}

// getCmd represents the get command
var cmdIloLogin = &cobra.Command{
	Use:   "login",
	Short: "Login to iLO: hpecli ilo login",
	RunE:  runLogin,
}

func runLogin(cmd *cobra.Command, args []string) error {

	if host == "" {
		return fmt.Errorf("must provide --host or -h")
	}
	if !strings.HasPrefix(host, "http") {
		host = fmt.Sprintf("http://%s", host)
	}

	if username == "" {
		return fmt.Errorf("must provide --username or -u")
	}

	if password == "" {
		// this really isn't secure to provide on the command line
		// need to provide a way to read from stdin
		return fmt.Errorf("must provide --password or -p")
	}

	fmt.Println(fmt.Sprintf("Attempting login with user: %v, at: %v", username, host))

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
	return fmt.Sprintf("hpecli_ilo_token_%s", host)
}
