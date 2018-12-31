package api

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ganlvtech/go-exportable-cookiejar"
	"github.com/pkg/errors"
)

func Md5Sum(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func RsaEncrypt(publicKey []byte, origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

func Timestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func StringToReader(str string) io.Reader {
	return bytes.NewReader([]byte(str))
}

func SaveCookieJar(jar *cookiejar.Jar, url1 string) (string, error) {
	url2, err := url.Parse(url1)
	if err != nil {
		return "", err
	}
	cookies := jar.Cookies(url2)
	data, err := json.Marshal(cookies)
	if err != nil {
		return "", err
	}
	return string(data[:]), nil
}

func LoadCookieJar(j *cookiejar.Jar, str string, url1 string) error {
	url2, err := url.Parse(url1)
	if err != nil {
		return err
	}
	cookies := make([]*http.Cookie, 0)
	err = json.Unmarshal([]byte(str), &cookies)
	if err != nil {
		return err
	}
	j.SetCookies(url2, cookies)
	return nil
}
