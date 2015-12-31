package evexmlapi

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// Params is a map for a request's query parameters
type Params map[string]string

// Resource is a stuct for a resource
type Resource struct {
	path           string
	cacheDuration  int64
	requiredParams []string
	optionalParams []string
	model
	Api
}

type Api struct {
	protocol string
	method   string
	baseURL  string
	port     int
}

var (
	xmlAPI = Api{
		protocol: "https",
		method:   "GET",
		baseURL:  "https://api.eveonline.com/",
		port:     443,
	}
)

type paramError struct {
	params Params
	resourcePath string
}

func (pe paramError) Error() string {
	return fmt.Sprintf("%q", pe.params) + " not valid for " + pe.resourcePath
}


// VerifyParams verifies the required params are available
func (r Resource) verifyParams(params Params) error {
	if (len(r.requiredParams) + len(r.optionalParams))  < len(params) {
		return paramError{params: params, resourcePath: r.path}
	}
	for _, v := range r.requiredParams {
		if _, ok := params[v]; !ok || params[v] == "" {
			return paramError{params: params, resourcePath: r.path}
		}
	}
	return nil
}

// ToJSON takes xml data and a format struct
// then converts it into a json byte array
func (r Resource) JSON(xmldata []byte) ([]byte, error) {
	f := r.format
	err := xml.Unmarshal(xmldata, f)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// NewSkillQueue is a constructor for the SkillQueue resource
func NewSkillQueue() Resource {
	return Resource{
		path:           "char/SkillQueue.xml.aspx",
		cacheDuration:  3600,
		requiredParams: []string{"keyID", "vCode", "characterID"},
		model:          model{format: &SkillQueueFormat{}},
		Api:            xmlAPI,
	}
}

// NewServerStatus is a constructor for the ServerStatus resource
func NewServerStatus() Resource {
	return Resource{
		path:           "server/ServerStatus.xml.aspx",
		cacheDuration:  180,
		requiredParams: []string{},
		model:          model{format: &ServerStatusFormat{}},
		Api:            xmlAPI,
	}
}

// NewAPIKeyInfo is a constructor for the APIKeyInfo resource
func NewAPIKeyInfo() Resource {
	return Resource{
		path:           "account/APIKeyInfo.xml.aspx",
		cacheDuration:  300,
		requiredParams: []string{"keyID", "vCode"},
		model:          model{format: &APIKeyInfoFormat{}},
		Api:            xmlAPI,
	}
}
