package main

import (
	"fmt"
	"github.com/youtube/vitess/go/cache"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type image []byte

func (img image) Size() int {
	return len(img)
}

var (
	imageCache          = cache.NewLRUCache(1024 * 1024 * 512) // max memory useage 512m
	hitCount    float64 = 0                                    // not thread safe, but enough for analytic
	missedCount float64 = 0
)

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

	// w.Header().Set("Cache-Control", "max-age=2592000")
	w.Header().Set("Expires", "Fri, 30 Oct 2998 14:19:41")
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
