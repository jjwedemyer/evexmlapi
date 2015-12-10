package evexmlapi

import (
	"io/ioutil"
	"log"
	"fmt"
	"net/http"
	"crypto/tls"	
	"net/url"
	"strings"
	
	"github.com/jovon/eve-xmlapi-go/models"
	"github.com/jovon/eve-xmlapi-go/cache"
	
)

const dateForm = "2006-01-02 15:04:05 UTC"
type Api struct {
	protocol string
	method string
	baseURL string
	port int	
}

type HttpRequest struct {
	userAgent   string
	api     Api
	params      map[string]string
	header http.Header
	cache cache.Cache
	client *http.Client
	respHandler responseHandler
}

func (hr HttpRequest) New() *HttpRequest {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}
	httpRequest := HttpRequest{cache: cache.MemoryCache{}.New(), 
								client: &http.Client{Transport: tr},
								respHandler: Handle(),}
	httpRequest.SetBaseURL("https://api.eveonline.com/")
	return &httpRequest
}

type responseHandler func(resp *http.Response) ([]byte, error)

func Handle() responseHandler {
	return func (resp *http.Response) ([]byte, error) {
		body, err := ioutil.ReadAll(resp.Body)	
		if err != nil {
			log.Fatalf("Reading body: %s", err)
			return []byte{}, err
		}
		statusCode := resp.StatusCode
		if statusCode != 200 {
			return []byte{}, fmt.Errorf("Status code(%d) returned; Message(%q)", statusCode, body)
		}
		return body, nil
	}
}

func (hr HttpRequest) AllParams() map[string]string {
	return hr.params
}

func (hr *HttpRequest) Param(p string) string {
	return hr.params[p]
}

func (hr *HttpRequest) setCache(c string) {
	
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
		return "eve-xmlapi-GO/0.0.0 (jvnpackard@gmail.com)"
	} else {
		return hr.userAgent
	}
}

func (hr *HttpRequest) SetUserAgent(ua string) {
	hr.userAgent = ua
}

func (hr HttpRequest) BaseURL() string {
	if hr.api.baseURL == "" {
		return "https://api.eveonline.com/"
	} else {
		return hr.api.baseURL
	}
}

func (hr *HttpRequest) SetBaseURL(bu string) {
	hr.api.baseURL = bu
}

// FetchXML retrieves the XML file from the cache or makes a request
func (hr HttpRequest) FetchXML(m models.Model) ([]byte, error) {
	fullPath := hr.url(m.Path())
	cacheValue := hr.cache.Read(fullPath)
	if cacheValue == nil {	
		results, err := hr.MakeRequest(fullPath)
		if err != nil {
			return nil, err
		}
		hr.cache.Write(fullPath, results, m.CacheDuration())
		return results, nil
	} 
	return cacheValue, nil
}

// MakeRequest sends the http request to get the xml
func (hr HttpRequest) MakeRequest(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", path, nil)	
	req.Header.Set("User-Agent", hr.UserAgent())
	resp, err := hr.client.Do(req)
	if err != nil {
		log.Fatalf("Do Request: %s", err)
		return []byte{}, err
	}		
	defer resp.Body.Close()
	
	return hr.respHandler(resp)
}

// url builds the full url path
func (hr *HttpRequest) url(path string) string {
	slash := "/"

	if strings.HasPrefix(path, slash) && strings.HasSuffix(hr.BaseURL(), slash) {
		strings.TrimLeft(path, slash)
	} else if !strings.HasPrefix(path, slash) && !strings.HasSuffix(hr.BaseURL(), slash) {
		path = slash + path
	}

	return hr.BaseURL() + path + "?" + hr.queryString()
}

func (hr *HttpRequest) queryString() string {
	queryParams := hr.AllParams()
	v := url.Values{}
	for key, value := range queryParams {
		v.Set(key, value)
	}
	return v.Encode()
}
