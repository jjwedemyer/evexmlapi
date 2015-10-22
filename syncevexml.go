package syncevexml

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/jovon/syncevexml/types"
)

type RequestStruct  types.RequestStruct


func (rs *RequestStruct) Fetch() string {
	var f RequestXML
	if localPath(rs.XmlPath) {
		f = LocalXMLRequest{Path: rs.XmlPath}
	} else {
		urlPath := buildURL(rs.XmlPath, rs.Params)
		f = HttpXMLRequest{
			Path: urlPath, 
			HttpClient: http.Client{},
			UserAgent: "SQO-GO/0.0.0 (jvnpackard@gmail.com)",
		}
	}
	data, err := f.FetchXML()
	if err != nil {
		fmt.Printf("%q", err)
	}
	result, _ := XMLToJSON(data, rs.XmlStruct)
	return result
}

func localPath(path string) bool {
	pathRegExp := regexp.MustCompile(`(\A[a-zA-Z]:\\?)|(\A\.\/)|(\A..\/)|(\A..\\)`)
	return pathRegExp.MatchString(path)
}

func XMLToJSON(xmlStr []byte, v interface{}) (string, error) {
	err := xml.Unmarshal(xmlStr, v)
	if err != nil {
		fmt.Println("Error unmarshalling from XML", err)
		return "", err
	}
	result, err := json.Marshal(v)
	if err != nil {
		fmt.Println("Error marshalling to JSON", err)
		return "", err
	}
	return string(result), nil
}

type LocalXMLRequest struct {
	Path string
}

type HttpXMLRequest struct {
	Path 		string
	HttpClient	http.Client
	UserAgent 	string
}

type RequestXML interface {
	FetchXML() ([]byte, error)
}

func (lxr LocalXMLRequest) FetchXML() ([]byte, error) {
	return ioutil.ReadFile(lxr.Path)
}

func (hxr HttpXMLRequest) FetchXML() ([]byte, error) {
	req, err := http.NewRequest("GET", hxr.Path, nil)
	req.Header.Set("User-Agent", hxr.UserAgent)
	client := &hxr.HttpClient 
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
		return make([]byte, 1), err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return make([]byte, 1), err
	}
	fmt.Println("Fetching XML")
	return body, nil
}

func buildURL(path string, queryParams map[string]string) string {
	baseURL := "https://api.testeveonline.com/"
	if strings.HasPrefix(path, "/") {
		strings.TrimLeft(path, "/")
	}
	v := url.Values{}
	for key, value := range queryParams {
		v.Set(key, value)
	}

	return fmt.Sprint(baseURL, path, "?", v.Encode())
}
