//go:build embed

package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"path"
	"strings"

	"github.com/gohutool/boot4go-util/http"
	"github.com/valyala/fasthttp"
)

// Embed the entire html directory.
//
//go:embed html
var embeddedHTML embed.FS

var (
	licenseJSOverride  string
	endpointJSOverride string
)

func setLicenseJS(li string) {
	if strings.TrimSpace(li) == "" {
		licenseJSOverride = ""
		return
	}
	licenseJSOverride = fmt.Sprintf("\n\t\t\t\tmyConfig.li=\"%v\";\n", li)
}

func setEndpointJS(endpoint string) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		endpointJSOverride = ""
		return
	}

	host, port := "unix", "2375"
	parts := strings.Split(endpoint, ":")
	if len(parts) >= 1 && strings.TrimSpace(parts[0]) != "" {
		host = strings.TrimSpace(parts[0])
	}
	if len(parts) >= 2 && strings.TrimSpace(parts[1]) != "" {
		port = strings.TrimSpace(parts[1])
	}

	endpointJSOverride = fmt.Sprintf("\n\t\t\t\tlocal_node.node_host = \"%v\";\n\t\t\t\tlocal_node.node_port = \"%v\";\n", host, port)
}

func newStaticHandler() fasthttp.RequestHandler {
	sub, err := fs.Sub(embeddedHTML, "html")
	if err != nil {
		panic(err)
	}

	pathNotFound := func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.SetContentType("application/json;charset=utf-8")
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Write([]byte(http.Result.Fail(fmt.Sprintf("Page Not Found, %v %v", string(ctx.Method()), string(ctx.RequestURI()))).Json()))
	}

	return func(ctx *fasthttp.RequestCtx) {
		p := string(ctx.Path())
		if p == "" || p == "/" {
			p = "/index.html"
		}
		if strings.HasSuffix(p, "/") {
			p += "index.html"
		}

		// Remove leading slash.
		p = strings.TrimPrefix(p, "/")

		// Dynamic overrides.
		switch p {
		case "static/public/js/cubeui.li.js":
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.Response.Header.SetContentType("application/javascript; charset=utf-8")
			ctx.Write([]byte(licenseJSOverride))
			return
		case "api/node.config.js":
			if endpointJSOverride != "" {
				ctx.SetStatusCode(fasthttp.StatusOK)
				ctx.Response.Header.SetContentType("application/javascript; charset=utf-8")
				ctx.Write([]byte(endpointJSOverride))
				return
			}
		}

		// Basic path traversal protection.
		clean := path.Clean("/" + p)
		if strings.Contains(clean, "..") {
			pathNotFound(ctx)
			return
		}
		p = strings.TrimPrefix(clean, "/")

		data, ct, ok := readEmbeddedFile(sub, p)
		if !ok {
			pathNotFound(ctx)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusOK)
		if ct != "" {
			ctx.Response.Header.SetContentType(ct)
		}
		ctx.Write(data)
	}
}

func readEmbeddedFile(fsys fs.FS, p string) ([]byte, string, bool) {
	// If it's a directory, try index files.
	if fi, err := fs.Stat(fsys, p); err == nil && fi.IsDir() {
		for _, idx := range []string{"index.html", "index.hml"} {
			b, ct, ok := readEmbeddedFile(fsys, path.Join(p, idx))
			if ok {
				return b, ct, true
			}
		}
		return nil, "", false
	}

	f, err := fsys.Open(p)
	if err != nil {
		return nil, "", false
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, "", false
	}

	ext := strings.ToLower(path.Ext(p))
	ct := mime.TypeByExtension(ext)
	if ct == "" {
		// Reasonable fallback for common text types.
		switch ext {
		case ".js":
			ct = "application/javascript; charset=utf-8"
		case ".css":
			ct = "text/css; charset=utf-8"
		case ".html", ".htm", ".hml":
			ct = "text/html; charset=utf-8"
		case ".json":
			ct = "application/json; charset=utf-8"
		case ".svg":
			ct = "image/svg+xml"
		}
	}

	return b, ct, true
}
