package httpsimple

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/grokify/gohttp"
)

var rxHTTPURL = regexp.MustCompile(`^(?i)https?://`)

type SimpleRequest struct {
	Method  string
	URL     string
	Query   map[string][]string
	Headers map[string][]string
	Body    interface{}
	IsJSON  bool
}

func (req *SimpleRequest) Inflate() {
	req.Method = strings.ToUpper(strings.TrimSpace(req.Method))
	if len(req.Method) == 0 {
		req.Method = http.MethodGet
	}
	if req.Headers == nil {
		req.Headers = map[string][]string{}
	}
	if req.IsJSON {
		if vals, ok := req.Headers[gohttp.HeaderContentType]; !ok {
			req.Headers[gohttp.HeaderContentType] =
				[]string{gohttp.ContentTypeAppJSONUtf8}
		} else {
			haveCT := false
			for _, hval := range vals {
				hval = strings.ToLower(strings.ToLower(hval))
				if strings.Index(hval, gohttp.ContentTypeAppJSON) == 0 {
					haveCT = true
					break
				}
			}
			if !haveCT {
				req.Headers[gohttp.HeaderContentType] = append(
					req.Headers[gohttp.HeaderContentType],
					gohttp.ContentTypeAppJSON)
			}
		}
	}
}

func (req *SimpleRequest) BodyBytes() ([]byte, error) {
	if req.Body == nil {
		return []byte{}, nil
	} else if reqBodyBytes, ok := req.Body.([]byte); ok {
		return reqBodyBytes, nil
	} else if reqBodyString, ok := req.Body.(string); ok {
		return []byte(reqBodyString), nil
	}
	return json.Marshal(req.Body)
}

// SimpleClient provides a simple interface to making HTTP requests
// using `net/http`.
type SimpleClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewSimpleClient(httpClient *http.Client, baseURL string) SimpleClient {
	return SimpleClient{HTTPClient: httpClient, BaseURL: baseURL}
}

func (sc *SimpleClient) Get(reqURL string) (*http.Response, error) {
	return sc.Do(SimpleRequest{Method: http.MethodGet, URL: reqURL})
}

func (sc *SimpleClient) Do(req SimpleRequest) (*http.Response, error) {
	req.Inflate()
	bodyBytes, err := req.BodyBytes()
	if err != nil {
		return nil, err
	}
	reqURL := strings.TrimSpace(req.URL)
	if len(sc.BaseURL) > 0 {
		if len(reqURL) == 0 {
			reqURL = sc.BaseURL
		} else if !rxHTTPURL.MatchString(reqURL) {
			reqURL = gohttp.JoinAbsolute(sc.BaseURL, reqURL)
		}
	}
	if len(req.Query) > 0 {
		goURL, err := gohttp.URLAddQueryString(reqURL, req.Query)
		if err != nil {
			return nil, err
		}
		reqURL = goURL.String()
	}
	if sc.HTTPClient == nil {
		sc.HTTPClient = &http.Client{}
	}
	return gohttp.DoJSONSimple(
		sc.HTTPClient, req.Method, reqURL, req.Headers, bodyBytes)
}

func (sc *SimpleClient) DoJSON(req SimpleRequest, resBody interface{}) ([]byte, *http.Response, error) {
	resp, err := sc.Do(req)
	if err != nil {
		return []byte{}, nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return bytes, resp, err
	}
	err = json.Unmarshal(bytes, resBody)
	return bytes, resp, err
}

func Do(req SimpleRequest) (*http.Response, error) {
	sc := SimpleClient{}
	return sc.Do(req)
}
