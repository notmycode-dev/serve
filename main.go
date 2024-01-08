package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

var (
	folder = flag.String("folder", "./", "Set the folder to serve files from (default is \"./\")")
	host   = flag.String("host", "0.0.0.0", "Set the hostname (default \"localhost\")")
	port   = flag.Int("port", 8080, "Set the port (default 8080)")
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Server is listening on %s...\n", addr)
	if err := fasthttp.ListenAndServe(addr, requestHandlerWrapper); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
