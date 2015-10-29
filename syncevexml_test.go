package syncevexml

import (
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
	skillQueue  = models.NewSkillQueue()
	dateNow     string
)

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
		pathCorrect, _ := regexp.MatchString(skillQueue.Path(), r.RequestURI)
		if !pathCorrect {
			t.Errorf("Wrong URL. Expected %q. Got %q", skillQueue.Path(), r.RequestURI)
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
	dat, _ := skillQueue.FromXML(&v)

	if skillQueue.GetCachedUntil(dat) == "" {
		t.Errorf("Error wrong response: %q", dat)
	}
}

func TestUpdateCache(t *testing.T) {
	duration, _ := time.ParseDuration("1h")
	dateNow = time.Now().Add(duration).Format(dateForm)

	r := fmt.Sprintf(
		`{"currentTime":"2009-04-18 15:19:43",
		"rows":[{"queuePosition":"1","typeID":"11441","level":"3","startSP":"7072","endSP":"40000","startTime":"2009-03-18 02:01:06","endTime":"2009-03-18 15:19:21"},
		{"queuePosition":"2","typeID":"20533","level":"4","startSP":"112000","endSP":"633542","startTime":"2009-03-18 15:19:21","endTime":"2009-03-30 03:16:14"}],
		"cachedUntil":%q}`,
		dateNow)
	cachedUntil := dateNow
	httpRequest.SetCachedUntil(cachedUntil)
	lastId := httpRequest.MergeCache([]byte(r), db, skillQueue)

	id, err := httpRequest.CheckCache(db, skillQueue)
	if err != nil {
		t.Error("Error checking cache: ", err)
	}
	if id == 0 {
		t.Error("CheckCache returned 0")
	}
	if id != lastId {
		t.Errorf("MergeCache = %q and CheckCache = %q returned different IDs", lastId, id)
	}
}

func TestIgnoreCache(t *testing.T) {

	r := fmt.Sprintf(
		`{"currentTime":"2009-04-18 15:19:43",
		"rows":[{"queuePosition":"1","typeID":"11441","level":"3","startSP":"7072","endSP":"40000","startTime":"2009-03-18 02:01:06","endTime":"2009-03-18 15:19:21"},
		{"queuePosition":"2","typeID":"20533","level":"4","startSP":"112000","endSP":"633542","startTime":"2009-03-18 15:19:21","endTime":"2009-03-30 03:16:14"}],
		"cachedUntil":%q}`,
		dateNow)
	cachedUntil := dateNow
	httpRequest.SetCachedUntil(cachedUntil)
	lastId := httpRequest.MergeCache([]byte(r), db, skillQueue)

	id, err := httpRequest.CheckCache(db, skillQueue)
	if err != nil {
		t.Error("Error checking cache: ", err)
	}
	if id == 0 {
		t.Error("CheckCache returned 0")
	}
	if id == lastId {
		t.Errorf("MergeCache = %q and CheckCache = %q returned the same id", lastId, id)
	}
}
