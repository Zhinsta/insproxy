package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func proxyHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Server", "insproxy")
	if req.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	picUrl, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(req.URL.String(), "/"))
	if err != nil {
		log.Println("proxy url is not valid: ", req.URL)
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	insUrl, err := url.Parse("http://" + string(picUrl))
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

	w.Header().Set("Expires", "Fri, 30 Oct 2998 14:19:41")
	w.Write(body)

	resp.Body.Close()
}
