package api

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/bitly/go-simplejson"
	"github.com/ganlvtech/go-exportable-cookiejar"
	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
)

// iOS
const APP_KEY = "27eb53fc9058f8c3"
const APP_SECRET = "c2ed53a74eeefe3cf99fbd01d8c9c375"

// Android
// const APP_KEY = "1d8b6e7d45233436"
// const APP_SECRET = "560c52ccd288fed045859ed18bffd973"
// 云视听 TV
// const APP_KEY = "4409e2ce8ffd12b8"
// const APP_SECRET = "59b43e04ad6965f34319062b478f83dd"

type BilibiliApiClient struct {
	Client        http.Client
	Username      string
	Password      string
	AccessToken   string
	RefreshToken  string
	BiliJct       string
	DanmakuConfig DanmakuConfig
}

type DanmakuConfig struct {
	Length int
	Color  int
	Mode   int
}

func NewBilibiliApiClient(debug bool) *BilibiliApiClient {
	b := new(BilibiliApiClient)
	b.Client = http.Client{}
	b.Client.Jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	if debug {
		proxyStr := "http://localhost:8888"
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			panic(err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		b.Client.Transport = transport
	}

	return b
}

func (b *BilibiliApiClient) GetBiliJct() (string, error) {
	if b.BiliJct != "" {
		return b.BiliJct, nil
	}
	u, _ := url.Parse("https://api.live.bilibili.com")
	cookies := b.Client.Jar.Cookies(u)
	for _, cookie := range cookies {
		if cookie.Name == "bili_jct" {
			b.BiliJct = cookie.Value
			return b.BiliJct, nil
		}
	}
	return "", errors.New("cannot find bili_jct in cookies")
}

func (b *BilibiliApiClient) SendLiveMessage(roomId int, content string) error {
	v := url.Values{}
	v.Set("color", "16777215")
	v.Set("fontsize", "25")
	v.Set("mode", "1")
	v.Set("msg", content)
	v.Set("rnd", Timestamp())
	v.Set("roomid", strconv.Itoa(roomId))
	biliJct, err := b.GetBiliJct()
	if err != nil {
		return err
	}
	v.Set("csrf_token", biliJct)
	v.Set("csrf", biliJct)
	resp, err := b.Client.PostForm("https://api.live.bilibili.com/msg/send", v)
	if err != nil {
		return err
	}
	j, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return err
	}
	code, err := j.Get("code").Int()
	if err != nil {
		return errors.WithMessage(err, "cannot get result code")
	}
	if code != 0 {
		message, _ := j.Get("message").String()
		return errors.Errorf("send live message failed: %s", message)
	}
	return nil
}

func (b *BilibiliApiClient) GetDanmakuConfig(roomId int, content string) error {
	resp, err := b.Client.Get("https://api.live.bilibili.com/userext/v1/danmuConf/getAll")
	if err != nil {
		return err
	}
	j, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return err
	}
	code, err := j.Get("code").Int()
	if err != nil {
		return errors.WithMessage(err, "cannot get result code")
	}
	if code != 0 {
		message, _ := j.Get("message").String()
		return errors.Errorf("get danmaku config failed: %s", message)
	}
	b.DanmakuConfig.Length, err = j.Get("data").Get("length").Int()
	if err != nil {
		return errors.WithMessage(err, "cannot get danmaku length")
	}
	b.DanmakuConfig.Color, err = j.Get("data").Get("color").Int()
	if err != nil {
		return errors.WithMessage(err, "cannot get danmaku color")
	}
	b.DanmakuConfig.Mode, err = j.Get("data").Get("mode").Int()
	if err != nil {
		return errors.WithMessage(err, "cannot get danmaku mode")
	}
	return nil
}

func (b *BilibiliApiClient) SaveCookie() ([]byte, error) {
	j, ok := b.Client.Jar.(*cookiejar.Jar)
	if !ok {
		return []byte{}, errors.New("cookie jar type assertion failed")
	}
	return j.JsonSerialize()
}

func (b *BilibiliApiClient) LoadCookies(data []byte) error {
	j, ok := b.Client.Jar.(*cookiejar.Jar)
	if !ok {
		return errors.New("cookie jar type assertion failed")
	}
	return j.JsonDeserialize(data)
}
