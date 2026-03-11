//go:build !embed

package main

import (
	"fmt"

	httputil "github.com/gohutool/boot4go-util/http"
	"github.com/valyala/fasthttp"
)

func newStaticHandler() fasthttp.RequestHandler {
	fs := &fasthttp.FS{
		Root:               "./html",
		IndexNames:         []string{"index.html", "index.hml"},
		GenerateIndexPages: true,
		Compress:           false,
		AcceptByteRange:    false,
		PathNotFound: func(ctx *fasthttp.RequestCtx) {
			ctx.Response.Header.SetContentType("application/json;charset=utf-8")
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write([]byte(httputil.Result.Fail(fmt.Sprintf("Page Not Found, %v %v", string(ctx.Method()), string(ctx.RequestURI()))).Json()))
		},
	}

	return fs.NewRequestHandler()
}
