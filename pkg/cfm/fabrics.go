package cfm

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getFabricsCommand() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "fabrics",
		Short: "Get Fabrics from HPE CFM: hpe cfm get fabrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getFabrics()
		},
	}

	return cmd
}

func getFabrics() error {
	host, token, err := hostAndToken()

	if err != nil {
		logrus.Debugf("unable to retrieve apiKey because of: %v", err)
		return fmt.Errorf("unable to retrieve the last login for CFM." +
			"Please login to CFM using: hpe cfm login")
	}

	cfmClient := newCFMClientFromAPIKey(host, token)
	logrus.Warningf("Using CFM: %s", host)

	// Get Fabrics will get the Fabrics of CFM of type *cfm.Fabrics
	fabrics, cfmError := cfmClient.GetFabrics()

	if cfmError != nil {
		logrus.Warningf("Could not login to CFM %s", cfmError.Result)
		return errors.New(cfmError.Result)
	}

	// The fabrics need to be marshalled to get the []byte format
	data, _ := json.Marshal(fabrics)
	logrus.Infof("%s", data)

	return nil
}
