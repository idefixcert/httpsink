package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":19091", "HTTP network address")
	flag.Parse()

	http.HandleFunc("/", rootHandler)

	log.Printf("Starting server on %s\n", *addr)

	err := http.ListenAndServe(*addr, logRequest(http.DefaultServeMux))
	if err != nil {
		log.Fatalln(err)
	}
}

func rootHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "<h1>Hello World</h1><div>Welcome to whereever you are</div>")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s %v\n", r.RemoteAddr, r.Method, r.URL, r.Header)
		defer func() {
			_ = r.Body.Close()
		}()
		body, err := ioutil.ReadAll(r.Body)
		log.Printf("%s %s\n", err, string(body))
		handler.ServeHTTP(w, r)
	})
}
