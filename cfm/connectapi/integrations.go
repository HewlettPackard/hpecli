package cfm

import (
	"encoding/json"

	"github.com/HewlettPackard/hpecli/cfm/cfmutils"
)

// IntegrationsSchema marshals and unmarshals CFM integration
type IntegrationsSchema struct {
	Name        string   `json:"name,omitempty"`
	Author      string   `json:"author,omitempty"`
	Version     string   `json:"version,omitempty"`
	Features    []string `json:"features,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
	Email       string   `json:"email,omitempty"`
	Description string   `json:"description,omitempty"`
}

// IntegrationSetSchema marshals and unmarshals CFM integration set
type IntegrationSetSchema struct {
	Description  string               `json:"description,omitempty"`
	IsSelected   bool                 `json:"is_selected,omitempty"`
	Integrations []IntegrationsSchema `json:"integrations,omitempty"`
	UUID         string               `json:"uuid,omitempty"`
	Name         string               `json:"name,omitempty"`
}

// Integrations marshals and unmarshals the CFM integrations
type Integrations struct {
	Count  int                  `json:"count,omitempty"`
	Result []IntegrationsSchema `json:"result,omitempty"`
}

// IntegrationSet marshals and unmarshals the Integration Sets associated with the CFM
type IntegrationSet struct {
	Count  int                  `json:"count"`
	Result IntegrationSetSchema `json:"result"`
}

// IntegrationSets marshals and unmarshals the Integration Sets associated with the CFM
type IntegrationSets struct {
	Count  int                    `json:"count"`
	Result []IntegrationSetSchema `json:"result"`
}

// GetIntegrations returns the associated integrations of CFM
func (cfmClient *CFMClient) GetIntegrations() (*Integrations, *cfmutils.ResponseObject) {
	url := "https://" + cfmClient.Host + "/api/v1/integrations"

	var integrations Integrations

	byteResponse, statusCode := cfmClient.RestCall("GET", url, "")

	if statusCode == 200 {
		err := json.Unmarshal(byteResponse, &integrations)
		cfmutils.CheckError(err)
		return &integrations, nil
	}

	errorObject := cfmutils.SetErrorObject(byteResponse, statusCode)
	return nil, errorObject
}

// GetIntegrationSets returns the associated integration sets of the CFM
func (cfmClient *CFMClient) GetIntegrationSets() (*IntegrationSets, *cfmutils.ResponseObject) {
	url := "https://" + cfmClient.Host + "/api/v1/integration_sets"

	var integrationSets IntegrationSets

	byteResponse, statusCode := cfmClient.RestCall("GET", url, "")

	if statusCode == 200 {
		err := json.Unmarshal(byteResponse, &integrationSets)
		cfmutils.CheckError(err)
		return &integrationSets, nil
	}

	errorObject := cfmutils.SetErrorObject(byteResponse, statusCode)
	return nil, errorObject
}
