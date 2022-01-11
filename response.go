package gohttp

import (
	"encoding/json"
	"time"
)

// ResponseInfo is a generic struct to handle response info.
type ResponseInfo struct {
	Name       string            `json:"name,omitempty"` // to distinguish from other requests
	Method     string            `json:"method,omitempty"`
	URL        string            `json:"url,omitempty"`
	StatusCode int               `json:"statusCode,omitempty"`
	Time       time.Time         `json:"time,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

// ToJSON returns ResponseInfo as a JSON byte array, embedding json.Marshal
// errors if encountered.
func (resIn *ResponseInfo) ToJSON() []byte {
	bytes, err := json.Marshal(resIn)
	if err != nil {
		resIn2 := ResponseInfo{StatusCode: 500, Body: err.Error()}
		bytes, _ := json.Marshal(resIn2)
		return bytes
	}
	return bytes
}
