// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/HewlettPackard/hpecli/pkg/greenlake/client"
	"github.com/HewlettPackard/hpecli/pkg/logger"
	"github.com/HewlettPackard/hpecli/pkg/store"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	host      string
	userID    string
	secretKey string
	tenantID  string
	grantType string
)

func init() {
	cmdGreenLakeLogin.Flags().StringVar(&host, "host", "", "greenlake ip address")
	cmdGreenLakeLogin.Flags().StringVarP(&userID, "userid", "u", "", "greenlake userID")
	cmdGreenLakeLogin.Flags().StringVarP(&secretKey, "secretkey", "s", "", "greenlake secretKey")
	cmdGreenLakeLogin.Flags().StringVarP(&tenantID, "tenantid", "t", "", "greenlake tenantID")
}

// getCmd represents the get command
var cmdGreenLakeLogin = &cobra.Command{
	Use:   "login",
	Short: "Login to greenlake: hpecli greenlake login",
	RunE:  runLogin,
}

func runLogin(cmd *cobra.Command, args []string) error {

	logger.Info("greenlake/login called")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	reader := bufio.NewReader(os.Stdin)
	if host == "" {
		fmt.Print("Enter host:")
		host, _ = reader.ReadString('\n')
		if host == "\n" {
			logger.Debug("Unable to login must provide --host or -h")
			return fmt.Errorf("must provide --host or -h")
		}
	}
	if !strings.HasPrefix(host, "http") {
		host = fmt.Sprintf("http://%s", host)
	}

	if userID == "" {
		fmt.Print("Enter userID:")
		userID, _ = reader.ReadString('\n')
		if userID == "\n" {
			logger.Debug("Unable to login must provide --userID or -u")
			return fmt.Errorf("must provide --userID or -u")
		}
	}
	if secretKey == "" {
		fmt.Print("Enter secretKey:")
		bytePassword, err := terminal.ReadPassword(0)
		if err == nil {
			//fmt.Println("\nPassword typed: " + string(bytePassword))
		}
		secretKey = string(bytePassword)
		if secretKey == "" {
			logger.Debug("Unable to login must provide --secretkey or -s")
			return fmt.Errorf("must provide --secretkey or -s")
		}
	}

	if tenantID == "" {
		fmt.Print("\nEnter tenantID: ")
		tenantID, _ = reader.ReadString('\n')
		if tenantID == "\n" {
			logger.Debug("Unable to login must provide --tenantID or -t")
			return fmt.Errorf("must provide --tenantID or -t")
		}
	}

	fmt.Println(fmt.Sprintf("Attempting login with user: %v, at: %v", userID, host))

	grantType = "client_credentials"
	client := client.NewGreenLakeClient(grantType, userID, secretKey, tenantID, host)
	body, err := client.GetToken()
	if err != nil {
		logger.Debug("unable to get the token with the supplied credentials: %v", err)
		return err
	}
	var result map[string]string

	json.Unmarshal(body, &result)

	token := string(result["access_token"])

	key1 := "hpecli_greenlake_token_" + host
	key2 := "hpecli_greenlake_tenantid_" + host

	db, err := store.Open()
	if err != nil {
		logger.Debug("unable to open keystore: %v", err)
		return fmt.Errorf("%v", err)
	}

	db.Put(key1, token)
	db.Put(key2, tenantID)
	db.Close()
	return nil

}

func key() string {
	return fmt.Sprintf("hpecli_greenlake_token_%s", host)
}
