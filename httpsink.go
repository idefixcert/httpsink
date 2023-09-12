package main

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"time"
)

func main() {
	addr := flag.String("addr", ":19091", "HTTP network address")
	sleep := flag.Duration("sleep", 0, "sleeptime for non proxy function")
	flag.Parse()

	http.HandleFunc("/", rootHandler(*sleep))

	log.Printf("Starting server on %s\n", *addr)

	err := http.ListenAndServe(*addr, logRequest(http.DefaultServeMux))
	if err != nil {
		log.Fatalln(err)
	}
}

func rootHandler(sleep time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(sleep)
		_, _ = fmt.Fprintf(w, "")
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
		message := fmt.Sprintf("Time:%v\nUrl:\n%sRequest:\n%s\nResponse Code: %d\nResponse:\n%s\n\n", time.Now(), r.URL.Path, string(x), rec.Code, rec.Body.String())
		if rec.Code >= 500 {
			color.Error.Block(message)
		} else if rec.Code >= 400 {
			color.Warn.Block(message)
		} else if rec.Code >= 200 {
			color.Info.Block(message)
		}
		// this copies the recorded response to the response writer
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		_, _ = rec.Body.WriteTo(w)
	})
}
