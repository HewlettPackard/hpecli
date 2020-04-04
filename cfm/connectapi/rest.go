package cfm

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/HewlettPackard/hpecli/cfm/cfmutils"

	log "github.com/HewlettPackard/hpecli/cfm/logging"
)

// RestCall makes the REST Call for SDK
func (cfmClient *CFMClient) RestCall(method string, url string, payload string) ([]byte, int) {

	var (
		request *http.Request
		err     error
	)

	if cfmClient.Host == "" {
		log.Error("Host not found. Host: " + cfmClient.Host)
		responseBody := `{"result": "Host is empty or not found"}`
		body, _ := json.Marshal(responseBody)
		log.Error("Host " + cfmClient.Host + " not found or is empty")
		return body, 404
	}

	if method == "GET" || method == "DELETE" {
		request, err = http.NewRequest(method, url, nil)
	} else {
		data := []byte(payload)
		request, err = http.NewRequest(method, url, bytes.NewBuffer(data))
	}
	cfmutils.SetAuthHeader(request, cfmClient.Token)
	cfmutils.CheckError(err)
	log.Info("Making a " + method + " call on " + url + " to the CFM " + cfmClient.Host)
	return cfmutils.GetResponse(request)
}
