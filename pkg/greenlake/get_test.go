//(C) Copyright 2019 Hewlett Packard Enterprise Development LP

package greenlake

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_runGet(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runGet(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("runGet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
