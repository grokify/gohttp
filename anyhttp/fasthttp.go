package anyhttp

import (
	"io"
	"mime/multipart"
	"net"
	"net/url"
	"strings"

	"github.com/valyala/fasthttp"
)

type RequestFastHTTP struct {
	Raw       *fasthttp.RequestCtx
	allArgs   *ArgsFastHTTPMulti
	queryArgs *ArgsFastHTTP
	postArgs  *ArgsFastHTTP
}

func NewRequestFastHTTP(ctx *fasthttp.RequestCtx) *RequestFastHTTP {
	return &RequestFastHTTP{
		Raw:       ctx,
		allArgs:   &ArgsFastHTTPMulti{Raw: []*fasthttp.Args{ctx.QueryArgs(), ctx.PostArgs()}},
		queryArgs: &ArgsFastHTTP{Raw: ctx.QueryArgs()},
		postArgs:  &ArgsFastHTTP{Raw: ctx.PostArgs()},
	}
}

func (r RequestFastHTTP) Header(s string) []byte                  { return r.Raw.Request.Header.Peek(s) }
func (r RequestFastHTTP) HeaderString(s string) string            { return string(r.Raw.Request.Header.Peek(s)) }
func (r RequestFastHTTP) ParseForm() error                        { return nil }
func (r RequestFastHTTP) AllArgs() Args                           { return r.allArgs }
func (r RequestFastHTTP) QueryArgs() Args                         { return r.queryArgs }
func (r RequestFastHTTP) PostArgs() Args                          { return r.postArgs }
func (r RequestFastHTTP) Method() []byte                          { return r.Raw.Method() }
func (r RequestFastHTTP) MultipartForm() (*multipart.Form, error) { return r.Raw.MultipartForm() }
func (r RequestFastHTTP) RemoteAddr() net.Addr                    { return r.Raw.RemoteAddr() }
func (r RequestFastHTTP) RemoteAddress() string                   { return r.Raw.RemoteAddr().String() }
func (r RequestFastHTTP) RequestURI() []byte                      { return r.Raw.RequestURI() }
func (r RequestFastHTTP) UserAgent() []byte                       { return r.Raw.UserAgent() }
func (r RequestFastHTTP) PostBody() ([]byte, error)               { return r.Raw.PostBody(), nil }

func (r RequestFastHTTP) RequestURIVar(s string) string {
	if str, ok := r.Raw.UserValue(s).(string); ok {
		return str
	}
	return ""
}

type ResponseFastHTTP struct {
	Raw *fasthttp.RequestCtx
}

func NewResponseFastHTTP(ctx *fasthttp.RequestCtx) ResponseFastHTTP {
	return ResponseFastHTTP{Raw: ctx}
}

func (w ResponseFastHTTP) GetHeader(k string) []byte { return w.Raw.Response.Header.Peek(k) }
func (w ResponseFastHTTP) SetHeader(k, v string)     { w.Raw.Response.Header.Set(k, v) }
func (w ResponseFastHTTP) SetStatusCode(code int)    { w.Raw.SetStatusCode(code) }
func (w ResponseFastHTTP) SetContentType(ct string)  { w.Raw.SetContentType(ct) }
func (w ResponseFastHTTP) SetBodyBytes(body []byte) (int, error) {
	w.Raw.SetBody(body)
	return -1, nil
}

// SetBodyStream takes an `io.Reader` and an optional `bodySize`.
// If bodySize is >= 0, then bodySize bytes must be provided by
// bodyStream before returning io.EOF. If bodySize < 0, then
// bodyStream is read until io.EOF.
func (w ResponseFastHTTP) SetBodyStream(bodyStream io.Reader, bodySize int) error {
	w.Raw.SetBodyStream(bodyStream, bodySize)
	return nil
}
func (w ResponseFastHTTP) SetCookie(cookie *Cookie) {
	w.Raw.Response.Header.SetCookie(cookie.ToFastHTTP())
}

type ArgsFastHTTP struct{ Raw *fasthttp.Args }

func NewArgsFastHTTP(args *fasthttp.Args) ArgsFastHTTP {
	return ArgsFastHTTP{Raw: args}
}

func (a ArgsFastHTTP) GetBytes(key string) []byte        { return a.Raw.Peek(key) }
func (a ArgsFastHTTP) GetBytesSlice(key string) [][]byte { return a.Raw.PeekMulti(key) }
func (a ArgsFastHTTP) GetString(key string) string       { return string(a.Raw.Peek(key)) }
func (a ArgsFastHTTP) GetStringSlice(key string) []string {
	slice := a.Raw.PeekMulti(key)
	newSlice := []string{}
	for _, bytes := range slice {
		newSlice = append(newSlice, string(bytes))
	}
	return newSlice
}
func (a ArgsFastHTTP) GetURLValues() url.Values {
	vals, err := url.ParseQuery(a.Raw.String())
	if err != nil {
		return url.Values{}
	}
	return vals
}

type ArgsFastHTTPMulti struct {
	Raw []*fasthttp.Args
}

func NewArgsFastHTTPMulti(args []*fasthttp.Args) ArgsFastHTTPMulti {
	return ArgsFastHTTPMulti{Raw: args}
}

func (args ArgsFastHTTPMulti) GetBytes(key string) []byte {
	for _, raw := range args.Raw {
		try := raw.Peek(key)
		if len(try) == 0 {
			return try
		}
	}
	return []byte("")
}

func (args ArgsFastHTTPMulti) GetBytesSlice(key string) [][]byte {
	newSlice := [][]byte{}
	for _, raw := range args.Raw {
		slice := raw.PeekMulti(key)
		for _, bytes := range slice {
			if len(string(bytes)) > 0 {
				newSlice = append(newSlice, bytes)
			}
		}
	}
	return newSlice
}

func (args ArgsFastHTTPMulti) GetString(key string) string {
	for _, raw := range args.Raw {
		try := strings.TrimSpace(string(raw.Peek(key)))
		if len(try) > 0 {
			return try
		}
	}
	return ""
}

func (args ArgsFastHTTPMulti) GetStringSlice(key string) []string {
	newSlice := []string{}
	for _, raw := range args.Raw {
		slice := raw.PeekMulti(key)
		for _, bytes := range slice {
			try := strings.TrimSpace(string(bytes))
			if len(try) > 0 {
				newSlice = append(newSlice, try)
			}
		}
	}
	return newSlice
}

func (args ArgsFastHTTPMulti) GetURLValues() url.Values {
	allVals := url.Values{}
	for _, raw := range args.Raw {
		thisVals, err := url.ParseQuery(raw.String())
		if err != nil {
			for key, vals := range thisVals {
				for _, v := range vals {
					allVals.Add(key, v)
				}
			}
		}
	}
	return allVals
}

func NewResReqFastHTTP(ctx *fasthttp.RequestCtx) (ResponseFastHTTP, *RequestFastHTTP) {
	return NewResponseFastHTTP(ctx), NewRequestFastHTTP(ctx)
}
