package evexmlapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func XMLServerRequest() *HttpRequest {
	httpRequest := NewRequest()
	httpRequest.overrideBaseURL = "https://api.testeveonline.com/"
	return httpRequest
}

var (
	httpRequest    = XMLServerRequest()
	testFileFolder = "./xml_examples/"
	emptyParams    = map[string][]string{}
	authParams     = map[string][]string{
		"characterID": []string{os.Getenv("characterID")},
		"keyID":       []string{os.Getenv("keyID")},
		"vCode":       []string{os.Getenv("vCode")},
	}
	errorMessage = "Error(%q) \nfetching:\n httpRequest %+v\n resource %+v"
)

func XMLHTTPTestServer(t *testing.T, params map[string][]string, rs Resource, testfile string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(testFileFolder + testfile)
		if err != nil {
			t.Error(err)
		}
		w.Header().Set("Content-Type", "application/xml")
		_ = r.ParseForm()
		for key, value := range params {
			if r.Form[key][0] != value[0] {
				w.WriteHeader(404)
				t.Logf("Key %q does not match, expected %q, got %q", key, value, r.Form[key][0])
			}
		}
		for _, value := range rs.requiredParams {
			if _, ok := r.Form[value]; !ok {
				t.Errorf("Param %q is required", value)
			}
		}
		if rs.path != r.URL.Path[1:] {
			w.WriteHeader(404)
			t.Logf("Wrong URL Path. Expected %q. Got %q", rs.path, r.URL.Path[1:])
		}

		if r.UserAgent() != httpRequest.UserAgent() {
			t.Logf("UserAgent does not match")
		}
		fmt.Fprintln(w, string(data))
	}))
}

func TestFetch_skillQueue(t *testing.T) {
	skillQueue := NewCharSkillQueue()
	params := authParams
	httpRequest.params = params
	httpRequest.SetFileCache("", "test-")

	httpRequest.SetResponseHandler(func(resp *http.Response) error {
		f := handler()
		return f(resp)
	})

	if testing.Short() {
		ts := XMLHTTPTestServer(t, params, skillQueue, "skillQueue.xml")
		defer ts.Close()
		httpRequest.overrideBaseURL = ts.URL
	}

	v, err := Fetch(skillQueue, httpRequest)
	if err != nil {
		t.Errorf(errorMessage, err, httpRequest, skillQueue)
	}
	if v == nil {
		t.Error("returned nil")
	}

	// Tests that an error is returned if no params are sent
	httpRequest.params = emptyParams
	v, err = Fetch(skillQueue, httpRequest)
	if err == nil {
		t.Errorf("Should be an Error, instead: %q", v)
	}
	// Tests that an error is returned if the wrong params are sent
	httpRequest.params = Params{"keyID": []string{""}}
	v, err = Fetch(skillQueue, httpRequest)
	if err == nil {
		t.Errorf("Should be an Error, instead: %q", err)
	}
}

func TestFetch_serverStatus(t *testing.T) {
	serverStatus := NewServerStatus()
	if testing.Short() {
		ts := XMLHTTPTestServer(t, emptyParams, serverStatus, "serverStatus.xml")
		defer ts.Close()
		httpRequest.overrideBaseURL = ts.URL
	}
	httpRequest.SetMemoryCache()
	httpRequest.params = emptyParams
	httpRequest.SetParser(func(data []byte, r Resource) ([]byte, error) {		
		return XMLtoJSON(data, r)
	})
	
	v, err := Fetch(serverStatus, httpRequest)
	if err != nil {
		t.Errorf(errorMessage, err, httpRequest, serverStatus)
		return
	}
	if v == nil {
		t.Error("returned nil")
		return
	}
	d := &ServerStatusFormat{}
	err = json.Unmarshal(v, d)
	if err != nil {
		t.Error(err)
		return
	}
	if d.CachedUntil == "" || d.OnlinePlayers == 0 {
		t.Error("Error with json unmarshal")
		return
	}
}

func TestFetch_apiKeyInfo(t *testing.T) {
	apiKeyInfo := NewAPIKeyInfo()
	params := Params{"keyID": authParams["keyID"], "vCode": authParams["vCode"]}
	if testing.Short() {
		ts := XMLHTTPTestServer(t, params, apiKeyInfo, "APIKeyInfo.xml")
		defer ts.Close()
		httpRequest.overrideBaseURL = ts.URL
	}
	httpRequest.params = params
	v, err := Fetch(apiKeyInfo, httpRequest)
	if err != nil {
		t.Errorf(errorMessage, err, httpRequest, apiKeyInfo)
		return
	}
	if v == nil {
		t.Error("returned nil")
		return
	}
	d := &APIKeyInfoFormat{}
	err = json.Unmarshal(v, d)
	if err != nil {
		t.Error(err)
		return
	}
	if d.CachedUntil == "" || len(d.Keys) == 0 {
		t.Error("Error with json unmarshal")
		return
	}
}
