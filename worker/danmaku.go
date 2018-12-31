package worker

import (
	"log"

	"github.com/ganlvtech/go-bilibili-helper/websocket"
)

func Danmaku() chan websocket.BilibiliPacket {
	c := make(chan websocket.BilibiliPacket)
	go func() {
		for packet := range c {
			switch packet.Op {
			case 3:
				popular := websocket.DecodePopular(packet.Data)
				log.Println("Received:", "Heartbeat Response.", "Popular:", popular)
			case 5:
				body := websocket.DecodeJSON(packet.Data)
				cmd, _ := body.Get("cmd").String()
				switch cmd {
				case "DANMU_MSG":
					text, _ := body.Get("info").GetIndex(1).String()
					uid, _ := body.Get("info").GetIndex(2).GetIndex(0).Int()
					username, _ := body.Get("info").GetIndex(2).GetIndex(1).String()
					medal, _ := body.Get("info").GetIndex(3).GetIndex(1).String()
					medalLevel, _ := body.Get("info").GetIndex(3).GetIndex(0).Int()
					userLevel, _ := body.Get("info").GetIndex(4).GetIndex(0).Int()
					log.Println("Received:", cmd, username, ":", text)
					log.Println("Received Danmaku User Info:", "[", medal, "|", medalLevel, "]", "[ UL", userLevel, "]", "uid =", uid)
				case "SEND_GIFT":
					username, _ := body.Get("data").Get("uname").String()
					action, _ := body.Get("data").Get("action").String()
					num, _ := body.Get("data").Get("num").Int()
					giftName, _ := body.Get("data").Get("giftName").String()
					price, _ := body.Get("data").Get("price").Int()
					totalCoin, _ := body.Get("data").Get("total_coin").Int()
					coinType, _ := body.Get("data").Get("coin_type").String()
					uid, _ := body.Get("data").Get("uid").Int()
					remain, _ := body.Get("data").Get("remain").Int()
					silver, _ := body.Get("data").Get("silver").Int()
					gold, _ := body.Get("data").Get("gold").Int()
					log.Println("Received:", cmd, username, action, num, giftName)
					log.Println("Received Gift:", num, "x", price, "=", totalCoin, coinType)
					log.Println("Received Gift User Info:", "uid =", uid, "remain =", remain, "silver =", silver, "gold =", gold)
				case "WELCOME":
					username, _ := body.Get("data").Get("uname").String()
					log.Println("Received:", cmd, username)
				default:
					pretty, _ := body.EncodePretty()
					log.Println("Received:", "Unknown Command", string(pretty))
				}
			default:
				log.Println("Unknown packet op")
			}
		}
	}()
	return c
}
