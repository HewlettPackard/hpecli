// (C) Copyright 2019 Hewlett Packard Enterprise Development LP.

package greenlake

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_runLogin(t *testing.T) {
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
			if err := runLogin(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("runLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_key(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := key(); got != tt.want {
				t.Errorf("key() = %v, want %v", got, tt.want)
			}
		})
	}
}
