package cfm

import (
	"encoding/json"
	"strings"

	"github.com/HewlettPackard/hpecli/cfm/cfmutils"

	log "github.com/HewlettPackard/hpecli/cfm/logging"
)

// UserSchema marshals and unmarshals the user schema associated with the CFM
type UserSchema struct {
	Username          string      `json:"username,omitempty"`
	TokenLifetime     int         `json:"token_lifetime"`
	UUID              string      `json:"uuid,omitempty"`
	AuthSourceUUID    string      `json:"auth_source_uuid,omitempty"`
	AuthSourceName    string      `json:"auth_source_name,omitempty"`
	DistinguishedName string      `json:"distinguished_name,omitempty"`
	Role              string      `json:"role,omitempty"`
	Immutable         bool        `json:"immutable,omitempty"`
	Preferences       interface{} `json:"preferences,omitempty"`
	Password          string      `json:"password,omitempty"`
}

// User unmarshals the user details
type User struct {
	Count  int        `json:"count,omitempty"`
	Result UserSchema `json:"result,omitempty"`
}

// Users unmarshals the users details
type Users struct {
	Count  int          `json:"count,omitempty"`
	Result []UserSchema `json:"result,omitempty"`
}

// GetUsers returns the details of all the users associated with the CFM
func (cfmClient *CFMClient) GetUsers(args ...string) (*Users, *cfmutils.ResponseObject) {
	url := "https://" + cfmClient.Host + "/api/v1/users"

	// The first query parameter would be the username
	if len(args) > 0 {
		username := args[0]
		username = strings.Replace(username, " ", "%20", -1)
		url += "/?username=" + username
	}

	var users Users

	byteResponse, statusCode := cfmClient.RestCall("GET", url, "")

	if statusCode == 200 {
		err := json.Unmarshal(byteResponse, &users)
		cfmutils.CheckError(err)
		log.Info("GetUsers successfully fetched the user information from CFM " + cfmClient.Host)
		return &users, nil
	}

	errorObject := cfmutils.SetErrorObject(byteResponse, statusCode)
	return nil, errorObject
}

// GetCurrentUser returns the details of the current user logged in into CFM
func (cfmClient *CFMClient) GetCurrentUser() (*User, *cfmutils.ResponseObject) {
	url := "https://" + cfmClient.Host + "/api/v1/users/current"

	var user User

	byteResponse, statusCode := cfmClient.RestCall("GET", url, "")

	if statusCode == 200 {
		err := json.Unmarshal(byteResponse, &user)
		cfmutils.CheckError(err)
		log.Info("GetCurrentUser successfully fetched the current user information")
		return &user, nil
	}

	errorObject := cfmutils.SetErrorObject(byteResponse, statusCode)
	return nil, errorObject
}

// GetUserByUUID (uuid string) returns the user details associated with the provided UUID
// Returns the response body of type UserSchema
func (cfmClient *CFMClient) GetUserByUUID(uuid string) (*User, *cfmutils.ResponseObject) {
	url := "https://" + cfmClient.Host + "/api/v1/users/" + uuid

	var user User

	byteResponse, statusCode := cfmClient.RestCall("GET", url, "")

	if statusCode == 200 {
		err := json.Unmarshal(byteResponse, &user)
		cfmutils.CheckError(err)
		log.Info("GetUserByUUID successfully fetched the user information of user " + uuid + " from CFM " + cfmClient.Host)
		return &user, nil
	}

	errorObject := cfmutils.SetErrorObject(byteResponse, statusCode)
	return nil, errorObject
}

// AddUser adds a new user to the CFM
// Returns the UUID of the user if success
func (cfmClient *CFMClient) AddUser(userData UserSchema) (*cfmutils.ResponseObject, *cfmutils.ResponseObject) {
	url := "https://" + cfmClient.Host + "/api/v1/users"

	var user cfmutils.ResponseObject
	payload, _ := json.Marshal(userData)

	byteResponse, statusCode := cfmClient.RestCall("POST", url, string(payload))

	if statusCode == 200 {
		err := json.Unmarshal(byteResponse, &user)
		cfmutils.CheckError(err)
		log.Info("Successfully added new user " + userData.Username + " to CFM " + cfmClient.Host)
		return &user, nil
	}

	errorObject := cfmutils.SetErrorObject(byteResponse, statusCode)
	return nil, errorObject
}

// DeleteUser deletes the user from CFM by UUID
// Returns nil if success
func (cfmClient *CFMClient) DeleteUser(uuid string) (*cfmutils.ResponseObject, *cfmutils.ResponseObject) {
	url := "https://" + cfmClient.Host + "/api/v1/users/" + uuid

	var result cfmutils.ResponseObject

	byteResponse, statusCode := cfmClient.RestCall("DELETE", url, "")

	if statusCode == 200 {
		err := json.Unmarshal(byteResponse, &result)
		cfmutils.CheckError(err)
		log.Info("Successfully deleted the user " + uuid + " from the CFM " + cfmClient.Host)
		return &result, nil
	}

	errorObject := cfmutils.SetErrorObject(byteResponse, statusCode)
	return nil, errorObject
}

// UpdateUser updates information of an existing user of CFM
// Returns nil if success
func (cfmClient *CFMClient) UpdateUser(uuid string, userData UserSchema) (*cfmutils.ResponseObject, *cfmutils.ResponseObject) {
	url := "https://" + cfmClient.Host + "/api/v1/users/" + uuid

	var user cfmutils.ResponseObject
	payload, _ := json.Marshal(userData)

	byteResponse, statusCode := cfmClient.RestCall("PUT", url, string(payload))

	if statusCode == 200 {
		err := json.Unmarshal(byteResponse, &user)
		cfmutils.CheckError(err)
		log.Info("Successfully updated the user " + uuid + " in the CFM " + cfmClient.Host)
		return &user, nil
	}

	errorObject := cfmutils.SetErrorObject(byteResponse, statusCode)
	return nil, errorObject
}

// ChangePassword updates the current password with the new password for the logged in user
// Returns a new token for the updated password if success
func (cfmClient *CFMClient) ChangePassword(currentPassword string, newPassword string) (*cfmutils.ResponseObject, *cfmutils.ResponseObject) {
	url := "https://" + cfmClient.Host + "/api/v1/users/current/change_password"

	var updatedPassword cfmutils.ResponseObject
	password := make(map[string]string)

	password["current_password"] = currentPassword
	password["new_password"] = newPassword

	payload, _ := json.Marshal(password)

	byteResponse, statusCode := cfmClient.RestCall("POST", url, string(payload))

	if statusCode == 200 {
		err := json.Unmarshal(byteResponse, &updatedPassword)
		cfmutils.CheckError(err)
		log.Info("Password successfully updated for the current user " + cfmClient.Username + " in the CFM " + cfmClient.Host)
		return &updatedPassword, nil
	}

	errorObject := cfmutils.SetErrorObject(byteResponse, statusCode)
	return nil, errorObject
}
