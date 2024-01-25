package main

import (
	"fmt"
	"html"
	"net/http"
	"os"

	"golang.org/x/time/rate"
)

func main() {
	http.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200"))
	})

	http.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	})

	http.HandleFunc("/404", http.NotFoundHandler().ServeHTTP)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		resp := []byte("<!DOCTYPE html>\n<html>\n")
		if r.Method == http.MethodGet {
			resp = append(resp, []byte("<em>Hello, world</em>\n<p>Query parameters:\n<ul>\n")...)
			for k, v := range r.URL.Query() {
				key := html.EscapeString(k)
				values := make([]string, len(v))
				for i, val := range v {
					values[i] = html.EscapeString(val)
				}
				resp = append(resp, []byte(fmt.Sprintf("<li>%s: %s</li>\n", key, values))...)
			}
			resp = append(resp, []byte("</ul>")...)
		} else {
			body := make([]byte, r.ContentLength)
			_, err := r.Body.Read(body)
			if err != nil {
				panic(err)
			}
			resp = append(resp, body...)
		}
		w.Write(resp)
	})

	http.HandleFunc("/authenticated", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			if username == os.Getenv("AUTH_USERNAME") && password == os.Getenv("AUTH_PASSWORD") {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("<!DOCTYPE html>\n<html>\nHello username!"))
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="localhost", charset="UTF-8"`)
		w.WriteHeader(http.StatusUnauthorized)
	})

	limiter := rate.NewLimiter(100, 30)
	http.HandleFunc("/limited", func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("200"))
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("429"))
		}

	})

	http.ListenAndServe(":8080", nil)
}
