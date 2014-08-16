package main

import (
	"compress/gzip"
	"fmt"
	"github.com/youtube/vitess/go/cache"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var (
	imageCache          = cache.NewLRUCache(1024 * 1024 * 512) // max memory useage 512m
	hitCount    float64 = 0                                    // not thread safe, but enough for analytic
	missedCount float64 = 0
)

type image []byte

func (img image) Size() int {
	return len(img)
}

func makeLogHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func makeGzipHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handler.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		handler.ServeHTTP(gzr, r)
	})
}

func proxyHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Server", "insproxy")
	if req.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	insUrl, err := url.Parse("http:/" + req.URL.String())
	if err != nil {
		log.Println("proxy url is not valid: ", req.URL)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if !isHostAllowed(insUrl.Host) {
		log.Println("proxy host is not allowed: ", req.URL)
		http.Error(w, "", http.StatusNotAcceptable)
		return
	}

	img, ok := imageCache.Get(insUrl.String())

	if ok {
		hitCount++
		w.Header().Set("Cache-Control", "max-age=2592000")
		w.Write(img.(image))
		return
	}
	missedCount++

	resp, err := http.Get(insUrl.String())
	if err != nil {
		log.Println("proxy failed: ", err)
		http.Error(w, "proxy failed", http.StatusBadGateway)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	imageCache.Set(insUrl.String(), image(body))

	w.Header().Set("Cache-Control", "max-age=2592000")
	w.Write(body)

	resp.Body.Close()
}

func statsHandler(w http.ResponseWriter, req *http.Request) {
	length, size, capacity, oldest := imageCache.Stats()
	hitRate := 0.0
	if missedCount == 0 && hitCount == 0 {
		hitRate = 0
	} else {
		hitRate = hitCount / (hitCount + missedCount)
	}
	fmt.Fprintf(w, `
    length: %v
    size: %v
    capacity: %v
    oldest: %v
    hit: %v
    missed: %v
    hit rate: %v
    `, length, size, capacity, oldest, hitCount, missedCount, hitRate)
}

func main() {
	http.HandleFunc("/", proxyHandler)
	http.HandleFunc("/stats", statsHandler)

	log.Println("About to listen 0.0.0.0:8080...")
	err := http.ListenAndServe(":8080", makeLogHandler(makeGzipHandler(http.DefaultServeMux)))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
