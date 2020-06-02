package main

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"time"
)

func main() {
	addr := flag.String("addr", ":19091", "HTTP network address")
	up := flag.String("up", "", "upstream address for proxy address, if not set, there is no upstream and the sink returns always with 200")
	sleep := flag.Duration("sleep", 0, "sleeptime for non proxy function")
	flag.Parse()

	if *up == "" {
		http.HandleFunc("/", rootHandler(*sleep))
	} else {
		http.HandleFunc("/", proxyHandler(*up))
	}

	log.Printf("Starting server on %s\n", *addr)

	err := http.ListenAndServe(*addr, logRequest(http.DefaultServeMux))
	if err != nil {
		log.Fatalln(err)
	}
}

func rootHandler(sleep time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(sleep)
		_, _ = fmt.Fprintf(w, "<h1>Hello World</h1><div>Welcome to whereever you are</div>")
	}
}

func proxyHandler(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// parse the url
		url, _ := url.Parse(target)

		// create the reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(url)

		// Update the headers to allow for SSL redirection
		req.URL.Host = url.Host
		req.URL.Scheme = url.Scheme
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = url.Host

		// Note that ServeHttp is non blocking and uses a go routine under the hood
		proxy.ServeHTTP(w, req)
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		x, err := httputil.DumpRequest(r, true)
		//defer func() {
		//	_ = r.Body.Close()
		//}()
		//body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		//log.Printf("%s %s\n", err, string(body))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, r)
		message := fmt.Sprintf("Request:\n%s\nResponse Code: %d\nResponse:\n%s\n\n", string(x), rec.Code, rec.Body.String())
		if rec.Code >= 500 {
			color.Error.Block(message)
		} else if rec.Code >= 404 {
			color.Warn.Block(message)
		} else if rec.Code >= 200 {
			color.Info.Block(message)
		}
		// this copies the recorded response to the response writer
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)
	})
}
