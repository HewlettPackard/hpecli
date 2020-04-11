package cfm

import (
	"encoding/json"
	"net/http"

	"github.com/HewlettPackard/hpecli/cfm/cfmutils"
)

// Authorization Unmarshals the response of session request
type Authorization struct {
	Count  int    `json:"count"`
	Result string `json:"result"`
}

// CFMClient is the return object for authentication that holds
// host, username, password and token information
type CFMClient struct {
	Host     string
	Username string
	Password string
	Token    string
}

// PasswordPolicyRules marshals and unmarhals the payload for Password Policy
type PasswordPolicyRules struct {
	MinimumUpperCase  int    `json:"minimum_upper_case,omitempty"`
	MinimumSpecial    int    `json:"minimum_special,omitempty"`
	MinimumLowerCase  int    `json:"minimum_lower_case,omitempty"`
	SpecialCharacters string `json:"special_characters,omitempty"`
	MinimumNumberic   int    `json:"minimum_numeric,omitempty"`
	MinimumLength     int    `json:"minimum_length,omitempty"`
	MaximumLength     int    `json:"maximum_length,omitempty"`
}

// PasswordPolicy marshals and unmarhals the payload for Password Policy
type PasswordPolicy struct {
	Count  int                 `json:"count,omitempty"`
	Result PasswordPolicyRules `json:"result,omitempty"`
}

// function setAuthHeader sets the authentication headers to get the authentication token from CFM
func setAuthHeader(request *http.Request, username, password string) {
	request.Header.Set("X-Auth-Username", username)
	request.Header.Set("X-Auth-Password", password)
}

// GetAuthToken generates and returns an authentication token for CFM API
func GetAuthToken(host, username, password string) (*CFMClient, *cfmutils.ResponseObject) {
	url := "https://" + host + "/api/v1/auth/token"

	var (
		auth Authorization
	)

	request, err := http.NewRequest("POST", url, nil)
	cfmutils.CheckError(err)
	cfmutils.SetHeader(request)
	setAuthHeader(request, username, password)

	byteResponse, statusCode := cfmutils.GetResponse(request)

	if statusCode == 200 {
		err = json.Unmarshal(byteResponse, &auth)
		cfmutils.CheckError(err)
		cfmClient := &CFMClient{host, username, password, auth.Result}
		return cfmClient, nil
	}

	errorObject := cfmutils.SetErrorObject(byteResponse, statusCode)
	return nil, errorObject
}

// GetPasswordPolicy returns the system password policy
func (cfmClient *CFMClient) GetPasswordPolicy() (*PasswordPolicy, *cfmutils.ResponseObject) {
	url := "https://" + cfmClient.Host + "/api/v1/auth/password_policy"

	var passwordPolicy PasswordPolicy

	byteResponse, statusCode := cfmClient.RestCall("GET", url, "")

	if statusCode == 200 {
		err := json.Unmarshal(byteResponse, &passwordPolicy)
		cfmutils.CheckError(err)
		return &passwordPolicy, nil
	}

	errorObject := cfmutils.SetErrorObject(byteResponse, statusCode)
	return nil, errorObject

}

// DeleteAuthToken returns a boolean value based on the deletion status of the token
func (cfmClient *CFMClient) DeleteAuthToken() (*cfmutils.ResponseObject, *cfmutils.ResponseObject) {
	url := "https://" + cfmClient.Host + "/api/v1/auth/token"

	var auth cfmutils.ResponseObject

	byteResponse, statusCode := cfmClient.RestCall("DELETE", url, "")

	if statusCode == 200 {
		err := json.Unmarshal(byteResponse, &auth)
		cfmutils.CheckError(err)
		return &auth, nil
	}

	errorObject := cfmutils.SetErrorObject(byteResponse, statusCode)
	return nil, errorObject
}
