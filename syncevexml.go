package syncevexml

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/jovon/syncevexml/db"
	"github.com/jovon/syncevexml/models"
	_ "github.com/lib/pq"
)

const dateForm = "2006-01-02 03:04:00"

type HttpRequest struct {
	userAgent string
	baseURL   string
	params    map[string]string
}

func testXMLServerRequest() *HttpRequest {
	httpRequest = HttpRequest{}
	httpRequest.SetBaseURL("https://api.testeveonline.com/")
	return &httpRequest
}

func (hxr HttpRequest) AllParams() map[string]string {
	return hxr.params
}

func (hxr *HttpRequest) Param(p string) string {
	return hxr.params[p]
}

func (hxr *HttpRequest) SetParams(params map[string]string) {
	if hxr.params == nil {
		hxr.params = map[string]string{}
	}
	for key, value := range params {
		hxr.params[key] = value
	}
}

func (hxr HttpRequest) UserAgent() string {
	if hxr.userAgent == "" {
		return "SQO-GO/0.0.0 (jvnpackard@gmail.com)"
	} else {
		return hxr.userAgent
	}
}

func (hxr *HttpRequest) SetUserAgent(ua string) {
	hxr.userAgent = ua
}

func (hxr HttpRequest) BaseURL() string {
	if hxr.baseURL == "" {
		return "https://api.testeveonline.com/"
	} else {
		return hxr.baseURL
	}
}

func (hxr *HttpRequest) SetBaseURL(bu string) {
	hxr.baseURL = bu
}

func (hxr *HttpRequest) XMLToJSON(xmlStr *[]byte, model models.Model) ([]byte, error) {
	v := model.DataFormat
	err := xml.Unmarshal(*xmlStr, &v)
	if err != nil {
		fmt.Println("Error unmarshalling from XML", err)
		return []byte{}, err
	}
	result, err := json.Marshal(v)
	if err != nil {
		fmt.Println("Error marshalling to JSON", err)
		return []byte{}, err
	}
	return result, nil
}

func (hxr HttpRequest) MergeCache(r []byte, db *data.DB, model models.Model) int64 {
	var lastId int64 = 0
	err := db.QueryRow(`Select merge_cache($1::integer, $2::varchar(25), $3::varchar(50), $4::jsonb)`,
		hxr.Param("keyID"), hxr.Param("characterID"), model.Path, string(r)).Scan(&lastId)

	if err != nil {
		log.Fatal("Fatal merge cache: ", err)
	}

	return lastId
}

func (hxr HttpRequest) CheckCache(db *data.DB, model models.Model) (int64, error) {
	var lastId int64 = 0
	err := db.QueryRow(`Select id from cache Where keyid = $1 and characterid = $2 and apipath = $3`,
		hxr.Param("keyID"), hxr.Param("characterID"), model.Path).Scan(&lastId)
	return lastId, err
}

func (hxr HttpRequest) Fetch(model models.Model) ([]byte, error) {
	fullPath := hxr.buildURL(model.Path)
	req, err := http.NewRequest("GET", fullPath, nil)

	req.Header.Set("User-Agent", hxr.UserAgent())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return []byte{}, err
	}
	return body, nil
}

func (hxr HttpRequest) buildURL(path string) string {
	queryParams := hxr.AllParams()
	if strings.HasPrefix(path, "/") {
		strings.TrimLeft(path, "/")
	}
	v := url.Values{}
	for key, value := range queryParams {
		v.Set(key, value)
	}

	return fmt.Sprint(hxr.BaseURL(), path, "?", v.Encode())
}
