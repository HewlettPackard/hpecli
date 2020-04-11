package cfm

import (
	"encoding/json"

	"github.com/HewlettPackard/hpecli/cfm/cfmutils"

	"github.com/sirupsen/logrus"
)

// Health unmarshals and marshals the fabric health
type Health struct {
	Status       string        `json:"status"`
	HealthIssues []interface{} `json:"health_issues"`
}

// FabricSchema marshals and unmarshals the fabric associated with CFM
type FabricSchema struct {
	Description            string `json:"description,omitempty"`
	ForeignManagerID       string `json:"foreign_manager_id,omitempty"`
	ForeignFabricState     string `json:"foreign_fabric_state,omitempty"`
	Name                   string `json:"name,omitempty"`
	Segmented              bool   `json:"segmented,omitempty"`
	Health                 Health `json:"health,omitempty"`
	FabricClass            string `json:"fabric_class,omitempty"`
	IsStable               bool   `json:"is_stable,omitempty"`
	ForeignManagementState string `json:"foreign_management_state,omitempty"`
	ForeignManagerURL      string `json:"foreign_manager_url,omitempty"`
	UUID                   string `json:"uuid,omitempty"`
}

// Fabric Unmarshals the Fabric Item associated with the CFM
type Fabric struct {
	Count  int          `json:"count,omitempty"`
	Result FabricSchema `json:"result,omitempty"`
}

// Fabrics Unmarshals the Fabric Collection associated with the CFM
type Fabrics struct {
	Count  int            `json:"count,omitempty"`
	Result []FabricSchema `json:"result,omitempty"`
}

// GetFabrics returns the collection of fabrics in the CFM
func (cfmClient *CFMClient) GetFabrics() (*Fabrics, *cfmutils.ResponseObject) {
	url := "https://" + cfmClient.Host + "/api/v1/fabrics"

	var fabrics Fabrics

	byteResponse, statusCode := cfmClient.RestCall("GET", url, "")

	if statusCode == 200 {
		err := json.Unmarshal(byteResponse, &fabrics)
		cfmutils.CheckError(err)
		return &fabrics, nil
	}

	errorObject := cfmutils.SetErrorObject(byteResponse, statusCode)
	logrus.Warnf("%s", errorObject.Result)
	return nil, errorObject
}

// GetFabricByUUID will return the fabric details associated with the uuid
// sent as parameter
func (cfmClient *CFMClient) GetFabricByUUID(uuid string) (*Fabric, *cfmutils.ResponseObject) {
	url := "https://" + cfmClient.Host + "/api/v1/fabrics/" + uuid

	var fabric Fabric

	byteResponse, statusCode := cfmClient.RestCall("GET", url, "")

	if statusCode == 200 {
		err := json.Unmarshal(byteResponse, &fabric)
		cfmutils.CheckError(err)
		return &fabric, nil
	}

	errorObject := cfmutils.SetErrorObject(byteResponse, statusCode)
	return nil, errorObject
}
