package main

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
)

func newProxyHandler(c *gin.Context) {
	encodedUrl := c.Params.ByName("url")

	picUrl, err := base64.StdEncoding.DecodeString(encodedUrl)
	if err != nil {
		c.String(406, "invalid url")
		return
	}

	insUrl, err := url.Parse("http://" + string(picUrl))
	if err != nil {
		c.String(406, "invalid url")
		return
	}

	if !isHostAllowed(insUrl.Host) {
		c.String(406, "invalid url")
		return
	}

	resp, err := http.Get(insUrl.String())
	if err != nil {
		c.String(502, "instagram error")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(500, "read error")
	}

	c.Data(200, resp.Header.Get("Content-Type"), body)
}

func main() {
	r := gin.Default()

	r.GET("/:url", newProxyHandler)

	r.Run(":8000")
}
