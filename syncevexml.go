package syncevexml

import (
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

const dateForm = "2006-01-02 15:04:05 UTC"

type HttpRequest struct {
	userAgent   string
	baseURL     string
	params      map[string]string
	cachedUntil string
}

func TestXMLServerRequest() *HttpRequest {
	httpRequest = HttpRequest{}
	httpRequest.SetBaseURL("https://api.testeveonline.com/")
	return &httpRequest
}

func (hr HttpRequest) AllParams() map[string]string {
	return hr.params
}

func (hr *HttpRequest) Param(p string) string {
	return hr.params[p]
}

func (hr *HttpRequest) SetParams(params map[string]string) {
	if hr.params == nil {
		hr.params = map[string]string{}
	}
	for key, value := range params {
		hr.params[key] = value
	}
}

func (hr HttpRequest) UserAgent() string {
	if hr.userAgent == "" {
		return "SQO-GO/0.0.0 (jvnpackard@gmail.com)"
	} else {
		return hr.userAgent
	}
}

func (hr *HttpRequest) SetCachedUntil(cu string) {
	hr.cachedUntil = cu
}

func (hr HttpRequest) CachedUntil() string {
	return hr.cachedUntil
}

func (hr *HttpRequest) SetUserAgent(ua string) {
	hr.userAgent = ua
}

func (hr HttpRequest) BaseURL() string {
	if hr.baseURL == "" {
		return "https://api.testeveonline.com/"
	} else {
		return hr.baseURL
	}
}

func (hr *HttpRequest) SetBaseURL(bu string) {
	hr.baseURL = bu
}

func (hr HttpRequest) MergeCache(r []byte, db *data.DB, model models.Model) int64 {
	var lastId int64 = 0
	// err := db.QueryRow(`Select merge_cache($1::integer, $2::varchar(25), $3::varchar(50), $4::jsonb, $5::timestamp)`,
	err := db.QueryRow(`Select insert_delete_cache($1::integer, $2::varchar(25), $3::varchar(50), $4::jsonb, $5::timestamp)`,
		hr.Param("keyID"), hr.Param("characterID"), model.Path, string(r), hr.CachedUntil()).Scan(&lastId)

	if err != nil {
		log.Fatal("Fatal merge cache: ", err)
	}

	return lastId
}

func (hr HttpRequest) CheckCache(db *data.DB, model models.Model) (int64, error) {
	var lastId int64 = 0
	err := db.QueryRow(`Select id from cache Where keyid = $1 and characterid = $2 and apipath = $3 and cachedUntil = $4`,
		hr.Param("keyID"), hr.Param("characterID"), model.Path, hr.CachedUntil()).Scan(&lastId)
	return lastId, err
}

func Poll(path string) string {
	resp, err := http.Head(path)
	if err != nil {
		log.Println("Error", path, err)
		return err.Error()
	}
	return resp.Status
}

func (hr HttpRequest) Fetch(model models.Model) ([]byte, error) {
	fullPath := hr.buildURL(model.Path)

	req, err := http.NewRequest("GET", fullPath, nil)

	req.Header.Set("User-Agent", hr.UserAgent())
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

func (hr HttpRequest) buildURL(path string) string {
	queryParams := hr.AllParams()
	if strings.HasPrefix(path, "/") {
		strings.TrimLeft(path, "/")
	}
	v := url.Values{}
	for key, value := range queryParams {
		v.Set(key, value)
	}

	return fmt.Sprint(hr.BaseURL(), path, "?", v.Encode())
}
