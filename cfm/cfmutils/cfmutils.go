package cfmutils

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/HewlettPackard/hpecli/cfm/logging"
)

// ResponseObject is the type for the error object sent if there is an error in API call
type ResponseObject struct {
	StatusCode int    `json:",omit_empty"`
	Result     string `json:"result,omitempty"`
}

// SetHeader will set the deafult headers for the CFM request
func SetHeader(request *http.Request) {
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Accept-encoding", "gzip, deflate, br")
	request.Header.Set("Accept-Language", "en-US,en;q=0.9,hi;q=0.8")
}

// SetAuthHeader adds Authorization token to the header along with
// other default headers
func SetAuthHeader(request *http.Request, token string) {
	SetHeader(request)
	request.Header.Set("Authorization", "Bearer "+token)
}

// SetErrorObject sets the error object when the status code is not 20x
func SetErrorObject(byteResponse []byte, statusCode int) *ResponseObject {
	var errorObject ResponseObject

	err := json.Unmarshal(byteResponse, &errorObject)
	CheckError(err)
	errorObject.StatusCode = statusCode
	log.Error(errorObject.Result)
	return &errorObject
}

// CheckError throws a panic if error is not nil
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// GetResponse parses makes the REST call, gets the response byte stream and
// returns the bytestream and the statuscode of the response
func GetResponse(request *http.Request) ([]byte, int) {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}

	client := &http.Client{Transport: transCfg}

	response, err := client.Do(request)
	client.CloseIdleConnections()
	CheckError(err)
	defer response.Body.Close()

	byteResponse, err := ioutil.ReadAll(response.Body)
	CheckError(err)
	log.Info("Request executed successfully")
	return byteResponse, response.StatusCode
}
