package anyhttp

import (
	"testing"
)

// TestInterface ensures following interface.
func TestInterface(t *testing.T) {
	nethttpReq := &RequestNetHTTP{}
	nethttpRes := ResponseNetHTTP{}

	MockRequest(nethttpReq)
	MockResponse(nethttpRes)
	MockHandler(nethttpRes, nethttpReq)

	fasthttpReq := RequestFastHTTP{}
	fasthttpRes := ResponseFastHTTP{}

	MockRequest(fasthttpReq)
	MockResponse(fasthttpRes)
	MockHandler(fasthttpRes, fasthttpReq)
}

func MockRequest(aReq Request)                {}
func MockResponse(aReq Response)              {}
func MockHandler(aRes Response, aReq Request) {}
