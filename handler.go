package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/valyala/fasthttp"
)

func listHandler(ctx *fasthttp.RequestCtx) {
	requestedPath := string(ctx.Path())
	fullPath := filepath.Join(*folder, requestedPath)
	ip := ctx.RemoteAddr().String()
	method := ctx.Method()
	log.Printf("Request from IP %s: %s %s", ip, method, ctx.Path())

	if fileInfos, err := ioutil.ReadDir(fullPath); err == nil {
		var htmlContent strings.Builder
		htmlContent.WriteString(fmt.Sprintf(`
			<head>
				<meta charset='UTF-8'>
				<meta name="viewport" content="width=device-width, initial-scale=1" />

				 <style>
				body {
					font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
					color: rgb(240, 240, 240);
					background-color: rgb(0, 0, 53);
					margin: 0;
					padding: 30px;
					-webkit-font-smoothing: antialiased;
				}
				
				ul li {
					list-style: none;
					justify-content: space-between;
					font-weight: 600;
					
				}
				
				ul a, span {
					color: rgb(254, 242, 255);
					text-decoration: none;
				}
				
				ul a:hover {
					text-decoration: underline;
				}
				
				div a {
					padding: 4px;
					border-radius: 10px;
				}
				
				
				</style>
				</head>
		`))

		htmlContent.WriteString(fmt.Sprintf("<h1>index of %s</h1>", requestedPath))
		htmlContent.WriteString("<div class='container'><ul>")

		for _, fileInfo := range fileInfos {
			name := fileInfo.Name()
			size := formatSize(fileInfo.Size())

			htmlContent.WriteString("<li>")
			if fileInfo.IsDir() {
				htmlContent.WriteString(fmt.Sprintf("<div><span>📁</span><a class='folder' href='%s/'>%s/</a></div>", filepath.Join(requestedPath, name), name))
			} else {
				htmlContent.WriteString(fmt.Sprintf("<div><span>📄</span><a class='file' href='%s'>%s</a> <span>%s</span></div>", filepath.Join(requestedPath, name), name, size))
			}
			htmlContent.WriteString("</li>")
		}
		htmlContent.WriteString("</ul></div>")

		ctx.Response.Header.Set("Content-Type", "text/html")
		ctx.Write([]byte(htmlContent.String()))
	} else {
		log.Printf("Error reading directory %s: %s", fullPath, err)
		ctx.Error(fmt.Sprintf("Error reading directory: %s", err), fasthttp.StatusInternalServerError)
		return
	}

	log.Printf("Request processed: %s %s", ctx.Method(), requestedPath)
}

func requestHandlerWrapper(ctx *fasthttp.RequestCtx) {
	listHandler(ctx)
}

func formatSize(size int64) string {
	const (
		B  = 1 << (10 * iota)
		KB = 1 << (10 * iota)
		MB = 1 << (10 * iota)
		GB = 1 << (10 * iota)
	)

	switch {
	case size < KB:
		return fmt.Sprintf("%d B", size)
	case size < MB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	case size < GB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	default:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	}
}
