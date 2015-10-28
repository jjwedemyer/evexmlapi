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

	"github.com/jovon/syncevexml/types"
	_ "github.com/lib/pq"
)

type HttpRequest types.HttpRequest

const (
	UserAgent string = "SQO-GO/0.0.0 (jvnpackard@gmail.com)"
	BaseURL   string = "https://api.testeveonline.com/"
)

func (hxr *HttpRequest) XMLToJSON(xmlStr *[]byte) ([]byte, error) {
	v := hxr.DataFormat
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

func (hxr *HttpRequest) MergeCache(r []byte, db *types.DB) int64 {
	var lastId int64 = 0
	err := db.QueryRow(`Select merge_cache($1::integer, $2::varchar(25), $3::varchar(50), $4::jsonb)`,
		hxr.Params["keyID"], hxr.Params["characterID"], hxr.Path, string(r)).Scan(&lastId)

	if err != nil {
		log.Fatal("Fatal merge cache in Store: ", err)
	}

	return lastId
}

func (hxr *HttpRequest) CheckCache(db *types.DB) (int64, error) {
	var lastId int64
	err := db.QueryRow("Select id from cache Where keyid = $1 and characterid = $2 and apipath = $3", hxr.Params["keyID"], hxr.Params["characterID"], hxr.Path).Scan(&lastId)
	return lastId, err
}

// func (dbr DBRequest) Fetch() (*sql.Rows, error) {
// 	db, err := sql.Open("postgres", dbr.ConnString)
// 	rows, err := db.Query(dbr.QueryString).Scan(&dbr.DataFormat)
// 	defer rows.Close()
// 	if err != nil {
// 		fmt.Printf("%q", err)
// 	}
// 	return rows, nil
// }

func (hxr HttpRequest) Fetch() ([]byte, error) {
	fullPath := hxr.buildURL()
	req, err := http.NewRequest("GET", fullPath, nil)
	req.Header.Set("User-Agent", UserAgent)
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

func (hxr HttpRequest) buildURL() string {
	var baseURL string
	if hxr.BaseURL != "" {
		baseURL = hxr.BaseURL
	} else {
		baseURL = BaseURL
	}

	path := hxr.Path
	queryParams := hxr.Params
	if strings.HasPrefix(path, "/") {
		strings.TrimLeft(path, "/")
	}
	v := url.Values{}
	for key, value := range queryParams {
		v.Set(key, value)
	}

	return fmt.Sprint(baseURL, path, "?", v.Encode())
}
