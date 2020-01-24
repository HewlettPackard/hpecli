//(C) Copyright 2019 Hewlett Packard Enterprise Development LP

package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/HewlettPackard/hpecli/pkg/logger"
)

const (
	contentType = "Content-Type"
	jsonType    = "application/json"
)

func TestNewGreenLakeClient(t *testing.T) {

	tests := []struct {
		name string
		args Client
		want *Client
	}{
		// TODO: Add more test cases.
		{
			name: "all inputs are provided for client",
			args: Client{
				GrantType:    "foo",
				ClientID:     "foo",
				ClientSecret: "foo",
				TenantID:     "foo",
				Host:         "foo",
			},
			want: &Client{
				GrantType:    "foo",
				ClientID:     "foo",
				ClientSecret: "foo",
				TenantID:     "foo",
				Host:         "foo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGreenLakeClient(tt.args.GrantType, tt.args.ClientID, tt.args.ClientSecret, tt.args.TenantID, tt.args.Host); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGreenLakeClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

// e.g. http.HandleFunc("/identity/v1/token", GetTokenHandler)
func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set(contentType, jsonType)
	jsonData := &Token{AccessToken: "abcdef",
		TokenType:       "Bearer",
		Expiry:          "date",
		ExpiresIn:       600,
		Scope:           "hpe-tenant",
		AccessTokenOnly: "false",
	}
	b, err := json.Marshal(jsonData)
	if err != nil {
		logger.Debug("unable to mashall the token data: %v", err)
		return
	}

	io.WriteString(w, string(b))
}
func TestClient_GetToken(t *testing.T) {
	tests := []struct {
		name    string
		fields  Client
		want    string
		wantErr bool
	}{
		// TODO: Add more test cases.
		{
			name: "token all inputs are provided",
			fields: Client{
				GrantType:    "foo",
				ClientID:     "foo",
				ClientSecret: "foo",
				TenantID:     "foo",
				Host:         "http://iam.intg.hpedevops.net",
			},
			want: `{"access_token":"abcdef","scope":"hpe-tenant","token_type":"Bearer","expiry":"date","expires_in":600,"accessTokenOnly":"false"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				GrantType:    tt.fields.GrantType,
				ClientID:     tt.fields.ClientID,
				ClientSecret: tt.fields.ClientSecret,
				TenantID:     tt.fields.TenantID,
				Host:         tt.fields.Host,
			}
			jsonValue, _ := json.Marshal(c)
			// Create a request to pass to our handler.
			req, err := http.NewRequest(http.MethodPost, "/identity/v1/token", bytes.NewBuffer(jsonValue))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetTokenHandler)
			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, req)

			got, err := c.GetToken()
			if err != nil {
				logger.Debug("unable to get the token with the supplied credentials: %v", err)
			}
			var result map[string]string

			json.Unmarshal(got, &result)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)

			}
			if rr.Body.String() != tt.want {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.want)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(result, tt.want) {
			// 	t.Errorf("Client.GetToken() = %v, want %v", result, tt.want)
			// }
		})
	}
}

// e.g. http.HandleFunc("/scim/v1/tenant/"+c.TenantID+"/"+tt.args.path", GetUsersHandler)
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set(contentType, jsonType)
	jsonData := &User{Active: true,
		DisplayName: "foo",
		UserName:    "foo",
		Name: Name{
			FamilyName: "foo",
			GivenName:  "foo",
		},
	}
	b, err := json.Marshal(jsonData)
	if err != nil {
		logger.Debug("unable to mashall the user data: %v", err)
		return
	}

	io.WriteString(w, string(b))
}

func TestClient_GetUsers(t *testing.T) {
	type args struct {
		path  string
		token string
	}
	tests := []struct {
		name    string
		fields  Client
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add more test cases.
		{
			name: "users all inputs are provided",
			fields: Client{
				GrantType:    "foo",
				ClientID:     "cli",
				ClientSecret: "key",
				TenantID:     "ten",
				Host:         "http://iam.intg.hpedevops.net",
			},
			args: args{
				path:  "Users",
				token: "abc",
			},
			want:    `{"active":true,"displayName":"foo","userName":"foo","name":{"familyName":"foo","givenName":"foo"}}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//response = nil
			c := &Client{
				GrantType:    tt.fields.GrantType,
				ClientID:     tt.fields.ClientID,
				ClientSecret: tt.fields.ClientSecret,
				TenantID:     tt.fields.TenantID,
				Host:         tt.fields.Host,
			}
			// Create a request to pass to our handler.
			req, err := http.NewRequest(http.MethodGet, c.Host+"/scim/v1/tenant/"+c.TenantID+"/"+tt.args.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tt.args.token)
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetUsersHandler)
			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, req)

			got, err := c.GetUsers(tt.args.path, tt.args.token)

			var result []User

			json.Unmarshal(got, &result)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)

			}
			if rr.Body.String() != tt.want {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.want)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Client.GetUsers() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestClient_doRequest(t *testing.T) {
	req1, _ := http.NewRequest(http.MethodPost, "http://iam.intg.hpedevops.net/identity/v1/token", nil)
	req2, _ := http.NewRequest(http.MethodGet, "http://iam.intg.hpedevops.net/scim/v1/tenant/ten/Users", nil)
	type args struct {
		req  *http.Request
		hand http.HandlerFunc
	}
	tests := []struct {
		name    string
		fields  Client
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Get Token request",
			fields: Client{
				GrantType:    "foo",
				ClientID:     "cli",
				ClientSecret: "key",
				TenantID:     "ten",
				Host:         "http://iam.intg.hpedevops.net",
			},
			args: args{
				req:  req1,
				hand: GetTokenHandler,
			},
		},
		{
			name: "Get Users request",
			fields: Client{
				GrantType:    "foo",
				ClientID:     "cli",
				ClientSecret: "key",
				TenantID:     "ten",
				Host:         "http://iam.intg.hpedevops.net",
			},
			args: args{
				req:  req2,
				hand: GetUsersHandler,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				GrantType:    tt.fields.GrantType,
				ClientID:     tt.fields.ClientID,
				ClientSecret: tt.fields.ClientSecret,
				TenantID:     tt.fields.TenantID,
				Host:         tt.fields.Host,
			}
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(tt.args.hand)
			handler.ServeHTTP(rr, tt.args.req)
			got, err := c.doRequest(tt.args.req)
			if err != nil {
				t.Errorf("Client.doRequest() error = %v, got %v", err, got)
				return
			}

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)

			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.doRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Client.doRequest() = %v, want %v", got, tt.want)
			// }
		})
	}
}
