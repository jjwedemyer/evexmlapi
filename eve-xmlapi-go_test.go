 package evexmlapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"	
	"os"	

	"github.com/jovon/eve-xmlapi-go/models"	
)

func XMLServerRequest() *HttpRequest {	
	httpRequest := HttpRequest{}.New()
	httpRequest.SetBaseURL("https://api.testeveonline.com/")
	return httpRequest
}

var (
	httpRequest = XMLServerRequest()
	skillQueue  = models.SkillQueue{}.New()
)

func TestFetch_http(t *testing.T) {
	params := map[string]string{
		"characterID": os.Getenv("characterID"),		
		"keyID":       os.Getenv("keyID"),
		"vCode":       os.Getenv("vCode"),
	}
	
	httpRequest.SetParams(params)

	if testing.Short() {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, _ := ioutil.ReadFile("./xml_examples/skillQueue.xml")
			w.Header().Set("Content-Type", "application/xml")
			_ = r.ParseForm()
			for key, value := range params {
				if r.Form[key][0] != value {
					w.WriteHeader(404)
					t.Logf("Key %q does not match, expected %q, got %q", key, value, r.Form[key][0])
				}
			}			
			if skillQueue.Path() != r.URL.Path[1:] {
				w.WriteHeader(404)
				t.Logf("Wrong URL Path. Expected %q. Got %q", skillQueue.Path(), r.URL.Path[1:])
			}
			
			if r.UserAgent() != httpRequest.UserAgent() {
				t.Logf("UserAgent does not match")
			}
			fmt.Fprintln(w, string(data))
		}))
		defer ts.Close()
		httpRequest.SetBaseURL(ts.URL)
	}

	v, err := httpRequest.FetchXML(skillQueue)
	if err != nil {
		t.Errorf("Error(%q) fetching: %q", err, httpRequest)
	}
	
	sk, err := skillQueue.ToStruct(v)
	if err != nil {
		t.Errorf("Error ToStruct: %q", err)
	}
	if len(sk.Rows) < 1 {
		t.Fail()
	}
	fmt.Printf("skillqueue: %q", sk)
	
	httpRequest.SetParams(map[string]string{"keyID": ""})
	v, err = httpRequest.FetchXML(skillQueue)
	if err == nil {
		t.Errorf("Should be an Error: %q", v)
	}	
}

