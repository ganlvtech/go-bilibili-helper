package api

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"

	"github.com/bitly/go-simplejson"
	"github.com/pkg/errors"
)

var PublicKey = ""
var Hash = ""

func GetRoomId(shortId int) int {
	resp, _ := http.Get("https://api.live.bilibili.com/room/v1/Room/room_init?id=" + strconv.Itoa(shortId))
	json2, _ := simplejson.NewFromReader(resp.Body)
	roomId, _ := json2.Get("data").Get("room_id").Int()
	return roomId
}

func SignPayload(payload map[string]string, accessToken string) url.Values {
	v := url.Values{}
	for key, value := range payload {
		if key != "sign" {
			v.Set(key, value)
		}
	}
	v.Set("access_key", accessToken)
	v.Set("actionKey", "appkey")
	v.Set("appkey", APP_KEY)
	v.Set("build", "8230")
	v.Set("device", "phone")
	v.Set("mobi_app", "iphone")
	v.Set("platform", "ios")
	v.Set("ts", Timestamp())
	v.Set("type", "json")
	// v.Encode() will sort params by key
	data := v.Encode()
	v.Set("sign", Md5Sum(data+APP_SECRET))
	return v
}

func GetPublicKey() (string, string, error) {
	if PublicKey != "" && Hash != "" {
		return PublicKey, Hash, nil
	}
	payload := make(map[string]string)
	resp, err := http.PostForm("https://passport.bilibili.com/api/oauth2/getKey", SignPayload(payload, ""))
	if err != nil {
		return "", "", err
	}
	j, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return "", "", err
	}
	code, err := j.Get("code").Int()
	if err != nil {
		return "", "", err
	}
	if code != 0 {
		message, err := j.Get("message").String()
		if err != nil {
			message = ""
		}
		return "", "", errors.New("get public key error: " + message)
	}
	PublicKey, err = j.Get("data").Get("key").String()
	if err != nil {
		return "", "", err
	}
	Hash, err = j.Get("data").Get("hash").String()
	if err != nil {
		return "", "", err
	}
	return PublicKey, Hash, nil
}

func EncryptPassword(password string) (string, error) {
	publicKey, hash, err := GetPublicKey()
	if err != nil {
		return "", err
	}
	crypt, err := RsaEncrypt([]byte(publicKey), []byte(hash+password))
	if err != nil {
		return "", err
	}
	passwordEncrypted := base64.StdEncoding.EncodeToString([]byte(crypt))
	return passwordEncrypted, nil
}
