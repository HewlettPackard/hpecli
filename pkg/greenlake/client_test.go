// (C) Copyright 2019 Hewlett Packard Enterprise Development LP

package greenlake

import (
	"reflect"
	"testing"
)

const (
	contentType            = "Content-Type"
	jsonType               = "application/json"
	host                   = "http://iam.intg.hpedevops.net"
	handlerWrongStatusCode = "handler returned wrong status code: got %v want %v"
)

func TestNewGreenLakeClient(t *testing.T) {
	tests := []struct {
		name string
		args Client
		want *Client
	}{
		{
			name: "all inputs are provided for client",
			args: Client{
				GrantType:    "foo",
				ClientID:     "foo",
				ClientSecret: "foo",
				TenantID:     "foo",
				Host:         host,
			},
			want: &Client{
				GrantType:    "foo",
				ClientID:     "foo",
				ClientSecret: "foo",
				TenantID:     "foo",
				Host:         host,
			},
		},
	}

	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGreenLakeClient(tt.args.GrantType, tt.args.ClientID,
				tt.args.ClientSecret, tt.args.TenantID, tt.args.Host); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGreenLakeClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
