package evexmlapi

import (
	"crypto/tls"
	// "io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/jovon/evexmlapi/cache"
)

const dateForm = "2006-01-02 15:04:05 UTC"

type HttpRequest struct {
	userAgent       string
	header          http.Header
	cache           cache.Cache
	client          *http.Client
	handlers        []ResponseStatusHandler
	parsers         []bodyParser
	params          Params
	fetch           Fetcher
	overrideBaseURL string
}

type ResponseStatusHandler func(resp *http.Response) error

type bodyParser func(b []byte, r Resource) ([]byte, error)

type Fetcher func(r Resource, hr *HttpRequest) ([]byte, error)

func XMLServerRequest() *HttpRequest {
	httpRequest := evexmlapi.NewRequest()
	httpRequest.overrideBaseURL = "https://api.testeveonline.com/"
	return httpRequest
}

// NewRequest constructs a new HttpRequest
func NewRequest() *HttpRequest {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}

	hR := HttpRequest{client: &http.Client{Transport: tr}}
	hR.cache = cache.NewMemoryCache()
	hR.params = Params{}
	hR.parsers = []bodyParser{rawParser()}
	hR.handlers = []ResponseStatusHandler{handler()}
	hR.fetch = fetchXML()
	return &hR
}

func handler() ResponseStatusHandler {
	return func(resp *http.Response) error {
		return nil
	}
}

func rawParser() bodyParser {
	return func(data []byte, r Resource) ([]byte, error) {
		return data, nil
	}
}

// SetToDefaultResponseHandler returns the response handler
// back to the default.
func (hr *HttpRequest) SetToDefaultResponseHandler() {
	hr.handlers = []ResponseStatusHandler{handler()}
}

// SetResponseHandler allows the user the change how the response
// is handled, mostly for handling the different response status codes.
func (hr *HttpRequest) SetResponseHandler(rh ResponseStatusHandler) {
	hr.handlers = append(hr.handlers, rh)
}

// SetParser allows the user to change how the responses body
// is parsed.
func (hr *HttpRequest) SetParser(parse bodyParser) {
	hr.parsers = append(hr.parsers, parse)
}

// UserAgent gets the userAgent
func (hr HttpRequest) UserAgent() string {
	if hr.userAgent == "" {
		return "eve-xmlapi-GO/0.0.0 (jvnpackard@gmail.com)"
	}
	return hr.userAgent
}

type cacheResults struct {
	key      string // used to find the results
	results  []byte // the value stored in cache
	duration int64  // time in seconds it should stay in cache
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
		read := <-readCh
		if read.err != nil {
			return nil, read.err
		}
		if read.results != nil {
			return read.results, nil
		}

		resp, err := hr.makeRequest(fullPath, r.method)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		results, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		for _, parse := range hr.parsers {
			results, err = parse(results, r)
			if err != nil {
				return nil, err
			}
		}

		writeCh := make(chan cacheResults)
		ca := cacheResults{key: fullPath, results: results, duration: r.cacheDuration}
		go hr.cacheResults(ca, writeCh)
		cached := <-writeCh
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
		ca.err = hr.cache.Store(ca.key, ca.results, ca.duration)
	}
	ch <- ca
}

// makeRequest sends the http request to get the xml
func (hr HttpRequest) makeRequest(path string, method string) (*http.Response, error) {
	req, err := http.NewRequest(method, path, nil)
	req.Header.Set("User-Agent", hr.UserAgent())
	resp, err := hr.client.Do(req)
	if err != nil {
		log.Fatalf("Do Request: %s", err)
		return nil, err
	}
	for _, h := range hr.handlers {
		err = h(resp)
		if err != nil {
			panic(err)
		}
	}
	return resp, nil
}

// url builds the full url path
func (hr HttpRequest) url(r Resource) (string, error) {
	slash := "/"
	path := r.path
	baseURL := r.protocol + r.baseURL
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
		for _, val := range value {
			v.Add(key, val)
		}
	}
	return v.Encode(), nil
}

func (hr *HttpRequest) SetFileCache(dir string, prefix string) {
	hr.cache = cache.NewFileCache(dir, prefix)
}

func (hr *HttpRequest) SetBoltCache(file string)  {
	hr.cache = cache.NewBoltCache(file, 0600, []byte("eve"), nil)
}

func (hr *HttpRequest) SetMemoryCache() {
	hr.cache = cache.NewMemoryCache()
}
