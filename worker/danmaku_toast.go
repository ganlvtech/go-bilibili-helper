package worker

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/ganlvtech/go-bilibili-helper/websocket"
)

type DanmakuToastInfo struct {
	Title string
	Text  string
}

func DanmakuToast() chan websocket.BilibiliPacket {
	q := make(chan DanmakuToastInfo, 10)
	go func() {
		for danmakuToastInfo := range q {
			err := exec.Command("SnoreToast.exe", "-t", danmakuToastInfo.Title, "-m", danmakuToastInfo.Text).Run()
			if err != nil {
				log.Println(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	c := make(chan websocket.BilibiliPacket)
	go func() {
		for packet := range c {
			if packet.Op == 5 {
				body := websocket.DecodeJSON(packet.Data)
				cmd, _ := body.Get("cmd").String()
				switch cmd {
				case "DANMU_MSG":
					text, _ := body.Get("info").GetIndex(1).String()
					username, _ := body.Get("info").GetIndex(2).GetIndex(1).String()
					medal, _ := body.Get("info").GetIndex(3).GetIndex(1).String()
					medalLevel, _ := body.Get("info").GetIndex(3).GetIndex(0).Int()
					userLevel, _ := body.Get("info").GetIndex(4).GetIndex(0).Int()
					title := fmt.Sprintf("[UL%d]%s", userLevel, username)
					if medal != "" {
						title = fmt.Sprintf("[%s|%d]%s", medal, medalLevel, title)
					}
					danmakuToastInfo := DanmakuToastInfo{
						title,
						text,
					}
					select {
					case q <- danmakuToastInfo:
					default:
					}
				case "SEND_GIFT":
					username, _ := body.Get("data").Get("uname").String()
					action, _ := body.Get("data").Get("action").String()
					num, _ := body.Get("data").Get("num").Int()
					giftName, _ := body.Get("data").Get("giftName").String()
					price, _ := body.Get("data").Get("price").Int()
					totalCoin, _ := body.Get("data").Get("total_coin").Int()
					coinType, _ := body.Get("data").Get("coin_type").String()
					remain, _ := body.Get("data").Get("remain").Int()
					silver, _ := body.Get("data").Get("silver").Int()
					gold, _ := body.Get("data").Get("gold").Int()
					free := ""
					if remain > 0 {
						free = "免费"
					}
					title := fmt.Sprintf("%s%s=%d个%s%s", username, action, num, free, giftName)
					text := fmt.Sprintf("%d*%d=%d %s", num, price, totalCoin, coinType)
					if remain > 0 {
						text += fmt.Sprintf("剩余%d个%s", remain, giftName)
					}
					if silver > 0 {
						text += fmt.Sprintf("剩余%d个银瓜子", silver)
					}
					if gold > 0 {
						text += fmt.Sprintf("剩余%d个金瓜子", gold)
					}
					danmakuToastInfo := DanmakuToastInfo{
						title,
						text,
					}
					select {
					case q <- danmakuToastInfo:
					default:
					}
				}
			}
		}
	}()
	return c
}
