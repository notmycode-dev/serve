package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/valyala/fasthttp"
)

var (
	folder = flag.String("folder", "./", "Set the folder to serve files from (default is \"./\")")
	host   = flag.String("host", "localhost", "Set the hostname (default \"localhost\")")
	port   = flag.Int("port", 8080, "Set the port (default 8080)")
)

func main() {
	flag.Parse()

	listHandler := func(ctx *fasthttp.RequestCtx) {
		requestedPath := string(ctx.Path())
		fullPath := filepath.Join(*folder, requestedPath)

		if fileInfo, err := os.Stat(fullPath); err == nil && fileInfo.IsDir() {
			fileInfos, err := os.ReadDir(fullPath)
			if err != nil {
				ctx.Error(fmt.Sprintf("Error reading directory: %s", err), fasthttp.StatusInternalServerError)
				return
			}

			var htmlContent strings.Builder
			htmlContent.WriteString(fmt.Sprintf(
				"<head>\n" +
					"<meta charset='UTF-8'>" +
					"    <link rel=\"stylesheet\" type=\"text/css\" href=\"/style.css\">\n" +
					"    <style>" + `
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
					
					ul a {
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
					
					
					` + "</style>" +
					"</head>\n",
			))

			htmlContent.WriteString(fmt.Sprintf("<h1>index of %s</h1>", requestedPath))
			htmlContent.WriteString("<div class='container'><ul>")

			for _, fileInfo := range fileInfos {
				name := fileInfo.Name()
				htmlContent.WriteString("<li>")
				if fileInfo.IsDir() {
					htmlContent.WriteString(fmt.Sprintf("<div><span>📁</span><a class='folder' href='%s/'>%s/</a></div>", filepath.Join(requestedPath, name), name))
				} else {
					htmlContent.WriteString(fmt.Sprintf("<div><span>📄</span><a class='file' href='%s'>%s</a></div>", filepath.Join(requestedPath, name), name))
				}
				htmlContent.WriteString("</li>")

			}
			htmlContent.WriteString("</ul></div>")

			ctx.Response.Header.Set("Content-Type", "text/html")
			ctx.Write([]byte(htmlContent.String()))
		} else {
			fileContent, err := ioutil.ReadFile(fullPath)
			if err != nil {
				ctx.Error(fmt.Sprintf("Error reading file: %s", err), fasthttp.StatusInternalServerError)
				return
			}

			contentType := mime.TypeByExtension(filepath.Ext(fullPath))
			if contentType == "" {
				contentType = "application/octet-stream"
			}
			ctx.Response.Header.Set("Content-Type", contentType)

			ctx.Write(fileContent)
		}
	}

	requestHandlerWrapper := func(ctx *fasthttp.RequestCtx) {
		listHandler(ctx)
	}

	addr := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("Server is listening on %s...\n", addr)
	if err := fasthttp.ListenAndServe(addr, requestHandlerWrapper); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
