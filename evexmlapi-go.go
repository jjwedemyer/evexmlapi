package evexmlapi

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/jovon/eve-xmlapi-go/cache"
)

const dateForm = "2006-01-02 15:04:05 UTC"

type HttpRequest struct {
	userAgent       string
	header          http.Header
	cache           cache.Cache
	client          *http.Client
	handle          ResponseHandler
	parse           BodyParser
	params          Params
	fetch           Fetcher
	overrideBaseURL string
}

type ResponseHandler func(resp *http.Response) (*http.Response, error)

type Body io.ReadCloser

type BodyParser func(b Body, r Resource) ([]byte, error)

type Fetcher func(r Resource, hr *HttpRequest) ([]byte, error)

func NewRequest() *HttpRequest {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}

	hR := HttpRequest{client: &http.Client{Transport: tr}}
	hR.SetCache("memory")
	hR.params = Params{}
	hR.parse = bodyParser()
	hR.handle = handler()
	hR.fetch = fetchXML()
	return &hR
}

func handler() ResponseHandler {
	return func(resp *http.Response) (*http.Response, error) {
		statusCode := resp.StatusCode
		if statusCode != 200 {
			defer resp.Body.Close()
			return nil, fmt.Errorf("Status code(%d) for request(%s) returned; Message(%q)", statusCode, resp.Request.URL.RawPath, resp.Body)
		}
		return resp, nil
	}
}

func bodyParser() BodyParser {
	return func(b Body, r Resource) ([]byte, error) {
		defer b.Close()
		return ioutil.ReadAll(b)
	}
}

// SetCache takes a string and an Options stuct with path and prefix
// properties for FileCache.
func (hr *HttpRequest) SetCache(c string) error {
	switch c {
	case "memory":
		hr.cache = cache.NewMemoryCache()
	case "file":
		hr.cache = cache.NewFileCache()
	default:
		return fmt.Errorf("setCache argument not recognized")
	}
	return nil
}

// SetToDefaultResponseHandler returns the response handler
// back to the default.
func (hr *HttpRequest) SetToDefaultResponseHandler() {
	hr.handle = handler()
}

func (hr *HttpRequest) SetResponseHandler(rh ResponseHandler) {
	hr.handle = rh
}

func (hr *HttpRequest) SetBodyParser(parse BodyParser) {
	hr.parse = parse
}

func (hr HttpRequest) UserAgent() string {
	if hr.userAgent == "" {
		return "eve-xmlapi-GO/0.0.0 (jvnpackard@gmail.com)"
	}
	return hr.userAgent
}

type cacheResults struct {
	key      string
	results  []byte
	duration int64
	err      error
}

// FetchXML retrieves the XML file from the cache or makes a request
func fetchXML() Fetcher {
	return func(r Resource, hr *HttpRequest) ([]byte, error) {
		err := r.verifyParams(hr.params)
		if err != nil {
			return nil, err
		}
		fullPath, err := hr.url(r)
		if err != nil {
			return nil, err
		}
		
		readCh := make(chan cacheResults)
		caResults := cacheResults{key: fullPath}
		go hr.checkCache(caResults, readCh)
		read := <- readCh
		if read.err != nil {
			return nil, read.err
		}
		if read.results != nil {
			return read.results, nil
		}
		
		resp, err := hr.makeRequest(fullPath)
		if err != nil {
			return nil, err
		}
		
		results, err := hr.parse(resp.Body, r)
		if err != nil {
			return nil, err
		}
		
		cacheCh := make(chan cacheResults)
		ca := cacheResults{key: fullPath, results: results, duration: r.cacheDuration}
		go hr.cacheResults(ca, cacheCh)
		cached := <-cacheCh
		if cached.err != nil {
			return nil, cached.err
		}
		return results, nil	
	}
}

func Fetch(r Resource, hr *HttpRequest) ([]byte, error) {
	return hr.fetch(r, hr)
}

func (hr HttpRequest) checkCache(ca cacheResults, ch chan cacheResults) {
	if hr.cache != nil {
		ca.results, ca.err = hr.cache.Read(ca.key)
	}
	ch <- ca
}

func (hr *HttpRequest) cacheResults(ca cacheResults, ch chan cacheResults) {
	if hr.cache != nil {
		err := hr.cache.Store(ca.key, ca.results, ca.duration)
		ca.err = err
	}
	ch <- ca
}

// makeRequest sends the http request to get the xml
func (hr HttpRequest) makeRequest(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", path, nil)
	req.Header.Set("User-Agent", hr.UserAgent())
	resp, err := hr.client.Do(req)
	if err != nil {
		log.Fatalf("Do Request: %s", err)
		return nil, err
	}
	return hr.handle(resp)
}

// url builds the full url path
func (hr HttpRequest) url(r Resource) (string, error) {
	slash := "/"
	path := r.path
	baseURL := r.baseURL
	if hr.overrideBaseURL != "" {
		baseURL = hr.overrideBaseURL
	}
	if strings.HasPrefix(path, slash) && strings.HasSuffix(baseURL, slash) {
		strings.TrimLeft(path, slash)
	} else if !strings.HasPrefix(path, slash) && !strings.HasSuffix(baseURL, slash) {
		path = slash + path
	}
	query, err := hr.queryString(r)
	if err != nil {
		return "", err
	}
	if query != "" {
		return baseURL + path + "?" + query, nil
	}
	return baseURL + path, nil
}

func (hr HttpRequest) queryString(r Resource) (string, error) {
	queryParams := hr.params
	
	v := url.Values{}
	for key, value := range queryParams {
		v.Set(key, value)
	}
	return v.Encode(), nil
}
