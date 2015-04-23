package main

import (
	"bytes"
	"crypto/des"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var ZHINSTA_CRYPT_KEY = os.Getenv("ZHINSTA_CRYPT_KEY")

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}

func urlDecrypt(encoded string) (decoded string, err error) {
	key := []byte(ZHINSTA_CRYPT_KEY)
	src, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	cipher, err := des.NewCipher(key)
	if err != nil {
		println("new des failed!")
		return "", err
	}

	out := make([]byte, len(src))
	dst := out

	bs := cipher.BlockSize()
	for len(src) > 0 {
		cipher.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return string(PKCS5UnPadding(out)), nil
}

func newProxyHandler(c *gin.Context) {
	if c.Request.Header.Get("Referer") != "http://www.zhinsta.com/" {
		c.String(406, "instersting ...")
		return
	}

	encodedUrl := c.Params.ByName("url")

	picUrl, err := urlDecrypt(encodedUrl)
	if err != nil {
		print(err)
		c.String(406, "invalid url in decrypt")
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
	defer resp.Body.Close()
	if err != nil {
		c.String(502, "instagram error")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(500, "read error")
	}

	// c.Writer.Header().Set("Cache-Control", "max-age=604800")
	c.Writer.Header().Set("Expires", "Fri, 15 May 2015 12:10:16 GMT")

	c.Data(200, resp.Header.Get("Content-Type"), body)
}

func main() {
	r := gin.Default()

	r.GET("/:url", newProxyHandler)

	r.Run(":8000")
}
