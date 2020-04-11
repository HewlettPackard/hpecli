package cfm

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/HewlettPackard/hpecli/cfm/cfmutils"
)

// RestCall makes the REST Call for SDK
func (cfmClient *CFMClient) RestCall(method string, url string, payload string) ([]byte, int) {

	var (
		request *http.Request
		err     error
	)

	if cfmClient.Host == "" {
		responseBody := `{"result": "Host is empty or not found"}`
		body, _ := json.Marshal(responseBody)
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
	return cfmutils.GetResponse(request)
}
