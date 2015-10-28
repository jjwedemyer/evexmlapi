package syncevexml

import (
	"bytes"
	// "database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jovon/syncevexml/types"
)

// type Env struct {
// 	db types.Cache
// }

func TestXMLToJSON(t *testing.T) {
	data, _ := ioutil.ReadFile("./xml_examples/skillQueue.xml")
	jsonResult := []byte(`{"currentTime":"2009-03-18 13:19:43","rows":[{"queuePosition":"1","typeID":"11441","level":"3","startSP":"7072","endSP":"40000","startTime":"2009-03-18 02:01:06","endTime":"2009-03-18 15:19:21"},{"queuePosition":"2","typeID":"20533","level":"4","startSP":"112000","endSP":"633542","startTime":"2009-03-18 15:19:21","endTime":"2009-03-30 03:16:14"}],"cachedUntil":"2009-03-18 13:34:43"}`)
	skillQueue := HttpRequest{
		UserAgent: "SQO-GO/0.0.0 (jvnpackard@gmail.com)",
		Params: map[string]string{
			"characterID": "95107920",
			"keyID":       "1469817",
			"vCode":       "7YIAEB9NpNoKKnn84kRrplBDNDQKtqzvIkqK8CsNxKFOVgtGOQJubdBQjuwCT9CN",
		},
		DataFormat: new(types.SkillQueueFormat),
	}
	result, err := skillQueue.XMLToJSON(&data)
	if !bytes.Equal(result, jsonResult) {
		t.Errorf("XMLToJSON(%q)\n EQUALS\n %q,\n WANT\n %q", data, result, jsonResult)
	}

	if err != nil {
		t.Errorf("%q", err)
	}
}

func TestFetch_http(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := ioutil.ReadFile("./xml_examples/skillQueue.xml")
		fmt.Fprintln(w, string(data))
	}))
	defer ts.Close()
	skillQueue := HttpRequest{
		Path:      "",
		UserAgent: "SQO-GO/0.0.0 (jvnpackard@gmail.com)",
		BaseURL:   ts.URL,
		Params: map[string]string{
			"characterID": "95107920",
			"keyID":       "1469817",
			"vCode":       "7YIAEB9NpNoKKnn84kRrplBDNDQKtqzvIkqK8CsNxKFOVgtGOQJubdBQjuwCT9CN",
		},
		DataFormat: new(types.SkillQueueFormat),
	}

	var dat map[string]interface{}
	v, err := skillQueue.Fetch()
	if err != nil {
		t.Errorf("Error fetching: %q", skillQueue)
	}
	j, _ := skillQueue.XMLToJSON(&v)
	if err = json.Unmarshal([]byte(j), &dat); err != nil {
		t.Errorf("Error Unmarshalling: %q\n; %q", v, err)
	}
	if dat["currentTime"] == "" {
		t.Errorf("Error wrong response: %q", dat)
	}
}

func TestCheckCache(t *testing.T) {
	db, err := types.NewDB()

	r := []byte(`{"currentTime":"2009-04-18 15:19:43","rows":[{"queuePosition":"1","typeID":"11441","level":"3","startSP":"7072","endSP":"40000","startTime":"2009-03-18 02:01:06","endTime":"2009-03-18 15:19:21"},{"queuePosition":"2","typeID":"20533","level":"4","startSP":"112000","endSP":"633542","startTime":"2009-03-18 15:19:21","endTime":"2009-03-30 03:16:14"}],"cachedUntil":"2009-03-18 13:34:43"}`)
	skillQueue := HttpRequest{
		UserAgent: "SQO-GO/0.0.0 (jvnpackard@gmail.com)",
		Path:      "char/SkillQueue.xml.aspx",
		Params: map[string]string{
			"characterID": "95107920",
			"keyID":       "1469817",
			"vCode":       "7YIAEB9NpNoKKnn84kRrplBDNDQKtqzvIkqK8CsNxKFOVgtGOQJubdBQjuwCT9CN",
		},
		DataFormat: new(types.SkillQueueFormat),
	}
	lastId := skillQueue.MergeCache(r, db)

	id, err := skillQueue.CheckCache(db)
	if err != nil {
		t.Errorf("Error checking cache: %q", err)
	}
	if id == 0 {
		t.Error("CheckCache returned 0")
	}
	if id != lastId {
		t.Error("Just Checking")
	}
}

// func TestRead(t *testing.T) {
// 	data, _ := ioutil.ReadFile("./xml_examples/skillQueue.xml")
// 	db, err := sql.Open(DbDialect, DbConnString)
// 	if err != nil {
// 		log.Fatal("Fatal Open DB:", err)
// 	}
// 	skillQueueHR := HttpRequest{
// 		Path: "char/SkillQueue.xml.aspx",
// 		Params: map[string]string{
// 			"characterID": "95107920",
// 			"keyID":       "1469817",
// 			"vCode":       "7YIAEB9NpNoKKnn84kRrplBDNDQKtqzvIkqK8CsNxKFOVgtGOQJubdBQjuwCT9CN",
// 		},
// 		DataFormat: new(types.SkillQueueFormat),
// 	}
// 	result, err := skillQueueHR.XMLToJSON(&data)
// 	if err != nil {
// 		t.Errorf("%q", err)
// 	}

// 	lastId := skillQueueHR.Store(result, db)
// 	if lastId == 0 {
// 		t.Error("Store did not work.")
// 	}
// }
