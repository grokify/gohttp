package gohttp

import (
	"bytes"
	"errors"
	"net/http"
	"strings"
)

func DoJSONSimple(client *http.Client, httpMethod, requrl string, headers map[string][]string, body []byte) (*http.Response, error) {
	requrl = strings.TrimSpace(requrl)
	if len(requrl) == 0 {
		return nil, errors.New("requrl is required but not present")
	}
	if client == nil {
		client = &http.Client{}
	}
	httpMethod = strings.TrimSpace(httpMethod)
	if httpMethod == "" {
		return nil, errors.New("httpMethod is required but not present")
	}
	var req *http.Request
	var err error

	if len(body) == 0 {
		req, err = http.NewRequest(httpMethod, requrl, nil)
	} else {
		req, err = http.NewRequest(httpMethod, requrl, bytes.NewBuffer(body))
	}
	if err != nil {
		return nil, err
	}
	for k, vals := range headers {
		k = strings.TrimSpace(k)
		kMatch := strings.ToLower(k)
		if kMatch == strings.ToLower(HeaderContentType) {
			continue
		}
		for _, v := range vals {
			req.Header.Set(k, v)
		}
	}
	if len(body) > 0 {
		req.Header.Set(HeaderContentType, ContentTypeAppJsonUtf8)
	}

	return client.Do(req)
}
