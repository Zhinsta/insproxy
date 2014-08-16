package main

import (
	"github.com/youtube/vitess/go/cache"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var imageCache = cache.NewLRUCache(1024)

type image []byte

func (img image) Size() int {
	return 1
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func proxyHandler(w http.ResponseWriter, req *http.Request) {
	println(imageCache.StatsJSON())
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
		log.Println("cache hit!")
		w.Write(img.(image))
		return
	}
	log.Println("cache missed!")

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

	w.Write(body)

	resp.Body.Close()
}

func main() {
	http.HandleFunc("/", proxyHandler)

	log.Println("About to listen 0.0.0.0:8080...")
	err := http.ListenAndServe(":8080", Log(http.DefaultServeMux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
