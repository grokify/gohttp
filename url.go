package gohttp

import (
	"net/url"
	"path"
	"regexp"
)

// JoinAbsolute performs a path.Join() while preserving two slashes after the scheme.
func JoinAbsolute(elem ...string) string {
	return regexp.MustCompile(`^([A-Za-z]+:/)`).ReplaceAllString(path.Join(elem...), "${1}/")
}

func URLAddQueryString(inputURL string, qry map[string][]string) (*url.URL, error) {
	goURL, err := url.Parse(inputURL)
	if err != nil {
		return nil, err
	}
	if len(qry) == 0 {
		return goURL, nil
	}
	allQS := goURL.Query()
	for k, vals := range qry {
		for _, val := range vals {
			allQS.Set(k, val)
		}
	}
	goURL.RawQuery = allQS.Encode()
	return goURL, nil
}
