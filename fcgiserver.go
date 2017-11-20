package main

import (
	"flag"
	"fmt"
	"github.com/jamesmcfadden/fcgiclient"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const FAST_CGI_HOST = "127.0.0.1"
const FAST_CGI_PORT = 9000

var root string
var staticMux *http.ServeMux

func serve(w http.ResponseWriter, r *http.Request) {
	abspath := root + r.URL.Path
	fcgi, _ := fcgiclient.New(FAST_CGI_HOST, FAST_CGI_PORT)

	// @todo - Use fast_cgi_params file - https://www.nginx.com/resources/wiki/start/topics/examples/phpfcgi/
	env := make(map[string]string)
	env["REQUEST_METHOD"] = r.Method
	env["SCRIPT_FILENAME"] = abspath

	f, _, _ := fcgi.Request(env, "")

	headers, body := parseFcgiResponse(string(f))

	for header, value := range headers {
		value = strings.TrimSpace(value)
		w.Header().Set(header, value)

		if header == "Status" {
			c, _ := strconv.Atoi(value[0:3])
			w.WriteHeader(c)
		}
	}

	io.WriteString(w, body)
}

func parseFcgiResponse(response string) (map[string]string, string) {
	chunks := strings.SplitN(response, "\r\n\r\n", 2)
	h := strings.Split(chunks[0], "\r\n")
	b := chunks[1]

	headers := make(map[string]string)

	for _, header := range h {
		headerParts := strings.Split(header, ":")
		headers[headerParts[0]] = headerParts[1]
	}

	return headers, b
}

func main() {

	listenPtr := flag.String("l", "localhost:8000", "The host/port to listen on")
	rootPtr := flag.String("r", "", "The web root")

	flag.Parse()

	root = *rootPtr

	fmt.Println("Listening on", *listenPtr)
	fmt.Println("Web root " + root)

	staticMux = http.NewServeMux()
	staticMux.Handle("/", http.FileServer(http.Dir(root)))

	mux := http.NewServeMux()
	mux.HandleFunc("/", serve)
	http.ListenAndServe(*listenPtr, mux)
}
