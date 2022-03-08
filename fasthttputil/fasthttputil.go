package fasthttputil

import (
	"strings"

	"github.com/valyala/fasthttp"
)

func GetReqQueryParam(ctx *fasthttp.RequestCtx, headerName string) string {
	return strings.TrimSpace(string(ctx.QueryArgs().Peek(headerName)))
}

func GetSplitReqQueryParam(ctx *fasthttp.RequestCtx, headerName, sep string) []string {
	return sliceTrimSpace(
		strings.Split(
			string(ctx.QueryArgs().Peek(headerName)), sep))
}

func sliceTrimSpace(s []string) []string {
	for i, v := range s {
		s[i] = strings.TrimSpace(v)
	}
	return s
}
