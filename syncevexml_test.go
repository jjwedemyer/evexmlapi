package syncevexml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/jovon/syncevexml/db"
	"github.com/jovon/syncevexml/models"
)

var (
	db          = data.NewDB()
	httpRequest = HttpRequest{}
	skillQueue  = models.SkillQueue()
)

func TestXMLToJSON(t *testing.T) {
	data, _ := ioutil.ReadFile("./xml_examples/skillQueue.xml")
	jsonResult := []byte(`{"currentTime":"2009-03-18 13:19:43","rows":[{"queuePosition":"1","typeID":"11441","level":"3","startSP":"7072","endSP":"40000","startTime":"2009-03-18 02:01:06","endTime":"2009-03-18 15:19:21"},{"queuePosition":"2","typeID":"20533","level":"4","startSP":"112000","endSP":"633542","startTime":"2009-03-18 15:19:21","endTime":"2009-03-30 03:16:14"}],"cachedUntil":"2009-03-18 13:34:43"}`)

	result, err := httpRequest.XMLToJSON(&data, skillQueue)
	if !bytes.Equal(result, jsonResult) {
		t.Errorf("XMLToJSON(%q)\n EQUALS\n %q,\n WANT\n %q", data, result, jsonResult)
	}

	if err != nil {
		t.Errorf("%q", err)
	}
}

func TestFetch_http(t *testing.T) {
	params := map[string]string{
		"characterID": "95107920",
		"keyID":       "1469817",
		"vCode":       "7YIAEB9NpNoKKnn84kRrplBDNDQKtqzvIkqK8CsNxKFOVgtGOQJubdBQjuwCT9CN",
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := ioutil.ReadFile("./xml_examples/skillQueue.xml")
		w.Header().Set("Content-Type", "application/xml")
		_ = r.ParseForm()
		for key, value := range params {
			if r.Form[key][0] != value {
				t.Errorf("Key %q does not match, expected %q, got %q", key, value, r.Form[key][0])
			}
		}
		pathCorrect, _ := regexp.MatchString(skillQueue.Path, r.RequestURI)
		if !pathCorrect {
			t.Errorf("Wrong URL. Expected %q. Got %q", skillQueue.Path, r.RequestURI)
		}
		if r.UserAgent() != httpRequest.UserAgent() {
			t.Error("UserAgent does not match")
		}
		fmt.Fprintln(w, string(data))
	}))
	defer ts.Close()

	httpRequest.SetParams(params)
	httpRequest.SetBaseURL(fmt.Sprint(ts.URL, "/"))

	v, err := httpRequest.Fetch(skillQueue)
	if err != nil {
		t.Errorf("Error fetching: %q", httpRequest)
	}
	j, _ := httpRequest.XMLToJSON(&v, skillQueue)
	dat := make(map[string]interface{})
	if err = json.Unmarshal([]byte(j), &dat); err != nil {
		t.Errorf("Error Unmarshalling: %q\n; %q", v, err)
	}
	if dat["currentTime"] == "" {
		t.Errorf("Error wrong response: %q", dat)
	}
}

func TestCheckCache(t *testing.T) {
	dateNow := time.Now().Format(dateForm)

	r := fmt.Sprintf(
		`{"currentTime":"2009-04-18 15:19:43",
		"rows":[{"queuePosition":"1","typeID":"11441","level":"3","startSP":"7072","endSP":"40000","startTime":"2009-03-18 02:01:06","endTime":"2009-03-18 15:19:21"},
		{"queuePosition":"2","typeID":"20533","level":"4","startSP":"112000","endSP":"633542","startTime":"2009-03-18 15:19:21","endTime":"2009-03-30 03:16:14"}],
		"cachedUntil":%q}`,
		dateNow)

	lastId := httpRequest.MergeCache([]byte(r), db, skillQueue)

	id, err := httpRequest.CheckCache(db, skillQueue)
	if err != nil {
		t.Errorf("Error checking cache: %q", err)
	}
	if id == 0 {
		t.Error("CheckCache returned 0")
	}
	if id != lastId {
		t.Error("MergeCache and CheckCache do not have the same result")
	}
}
