package main

import (
	"log"
	"net/http"
	"net/url"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func proxyHandler(w http.ResponseWriter, req *http.Request) {
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

	resp, err := http.Get(insUrl.String())
	if err != nil {
		log.Println("proxy failed: ", err)
		http.Error(w, "proxy failed", http.StatusBadGateway)
		return
	}
	log.Println(resp)
	resp.Write(w)
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
