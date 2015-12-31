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
	emptyParams    = map[string]string{}
	authParams     = map[string]string{
		"characterID": os.Getenv("characterID"),
		"keyID":       os.Getenv("keyID"),
		"vCode":       os.Getenv("vCode"),
	}
)

func XMLHTTPTestServer(t *testing.T, params map[string]string, rs Resource, testfile string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(testFileFolder + testfile)
		if err != nil {
			t.Error(err)
		}
		w.Header().Set("Content-Type", "application/xml")
		_ = r.ParseForm()
		for key, value := range params {
			if r.Form[key][0] != value {
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

func TestFetchXML_skillQueue(t *testing.T) {
	skillQueue := NewSkillQueue()
	params := authParams
	httpRequest.params = params

	httpRequest.SetResponseHandler(func(resp *http.Response) (*http.Response, error) {
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
		t.Errorf("Error(%q) fetching: %q", err, httpRequest)
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
	httpRequest.params = Params{"keyID": ""}
	v, err = Fetch(skillQueue, httpRequest)
	if err == nil {
		t.Errorf("Should be an Error, instead: %q", err)
	}
}

func TestFetchXML_serverStatus(t *testing.T) {
	serverStatus := NewServerStatus()
	httpRequest.SetCache("file")
	if testing.Short() {
		ts := XMLHTTPTestServer(t, emptyParams, serverStatus, "serverStatus.xml")
		defer ts.Close()
		httpRequest.overrideBaseURL = ts.URL
	}
	httpRequest.params = emptyParams
	httpRequest.SetBodyParser(func(b Body, r Resource) ([]byte, error) {
		defer b.Close()
		raw, err := ioutil.ReadAll(b)
		if err != nil {
			return raw, err
		}
		return r.JSON(raw)
	})
	v, err := Fetch(serverStatus, httpRequest)
	if err != nil {
		t.Errorf("Error(%q) fetching: %q", err, httpRequest)
	}
	if v == nil {
		t.Error("returned nil")
	}
	d := &ServerStatusFormat{}
	err = json.Unmarshal(v, d)
	if err != nil {
		t.Error(err)
	}
	if d.CachedUntil == "" || d.OnlinePlayers == 0 {
		t.Error("Error with json unmarshal")
	}
}

func TestFetchXML_apiKeyInfo(t *testing.T) {
	apiKeyInfo := NewAPIKeyInfo()
	httpRequest.SetCache("file")
	params := Params{"keyID": authParams["keyID"], "vCode": authParams["vCode"]}
	if testing.Short() {
		ts := XMLHTTPTestServer(t, params, apiKeyInfo, "APIKeyInfo.xml")
		defer ts.Close()
		httpRequest.overrideBaseURL = ts.URL
	}
	httpRequest.params = params
	v, err := Fetch(apiKeyInfo, httpRequest)
	if err != nil {
		t.Errorf("Error(%q) fetching: %q", err, httpRequest)
	}
	if v == nil {
		t.Error("returned nil")
	}
	d := &APIKeyInfoFormat{}
	err = json.Unmarshal(v, d)
	if err != nil {
		t.Error(err)
	}
	if d.CachedUntil == "" || len(d.Keys) == 0 {
		t.Error("Error with json unmarshal")
	}
}
