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

func main() {

	skillQueue := RequestStruct{
		// xmlPath: "char/SkillQueue.xml.aspx",
		XmlPath:   "./xml_examples/skillQueue.xml",
		XmlStruct: new(types.SkillQueueXML),
		Params: map[string]string{
			"characterID": "95107920",
			"keyID":       "1469817",
			"vCode":       "7YIAEB9NpNoKKnn84kRrplBDNDQKtqzvIkqK8CsNxKFOVgtGOQJubdBQjuwCT9CN",
		},
	}
	v := skillQueue.fetch()

	fmt.Printf("%s\n", v)
}

type LocalXMLPath struct {
	path string
}

type HttpXMLPath struct {
	path string
}

type XMLPath interface {
	fetchXML() ([]byte, error)
}

func (rs *RequestStruct) fetch() string {
	var f XMLPath
	if localPath(rs.XmlPath) {
		f = LocalXMLPath{path: rs.XmlPath}
	} else {
		urlPath := buildURL(rs.XmlPath, rs.Params)
		f = HttpXMLPath{path: urlPath}
	}
	data, err := f.fetchXML()
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

func (lxp LocalXMLPath) fetchXML() ([]byte, error) {
	return ioutil.ReadFile(lxp.path)
}

func (hxp HttpXMLPath) fetchXML() ([]byte, error) {
	req, err := http.NewRequest("GET", hxp.path, nil)
	req.Header.Set("User-Agent", "SQO-GO/0.0.0 (jvnpackard@gmail.com)")
	client := &http.Client{}
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
