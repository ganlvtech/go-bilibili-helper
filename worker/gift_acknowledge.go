package worker

import (
	"log"
	"strconv"
	"time"

	"github.com/ganlvtech/go-bilibili-helper/api"
	"github.com/ganlvtech/go-bilibili-helper/websocket"
)

func GiftAcknowledge(roomId int, b *api.BilibiliApiClient) chan websocket.BilibiliPacket {
	c := make(chan websocket.BilibiliPacket)
	q := make(chan string, 10)
	go func() {
		for content := range q {
			err := b.SendLiveMessage(roomId, content)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		for packet := range c {
			if packet.Op == 5 {
				body := websocket.DecodeJSON(packet.Data)
				cmd, _ := body.Get("cmd").String()
				if cmd == "SEND_GIFT" {
					username, _ := body.Get("data").Get("uname").String()
					action, _ := body.Get("data").Get("action").String()
					num, _ := body.Get("data").Get("num").Int()
					giftName, _ := body.Get("data").Get("giftName").String()
					content := "感谢" + username + action + "的" + strconv.Itoa(num) + "个" + giftName
					select {
					case q <- content:
					default:
					}
				}
			}
		}
	}()
	return c
}
