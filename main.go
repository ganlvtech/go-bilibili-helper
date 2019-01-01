package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ganlvtech/go-bilibili-helper/api"
	"github.com/ganlvtech/go-bilibili-helper/worker"
)

func main() {
	var err error

	// 参数
	// go-bilibili-helper.exe 23058 config.json
	if len(os.Args) < 2 {
		log.Println("没有指定房间号")
		return
	}
	if len(os.Args) < 3 {
		log.Println("没有指定配置文件")
		return
	}
	shortId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Println(err)
		return
	}
	configFile := os.Args[2]

	// 用户登录
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println(err)
	}
	config, err := LoadConfig(bytes)
	if err != nil {
		log.Println(err)
	}
	b := api.NewBilibiliApiClient(true)
	err = b.Login(config.Username, config.Password, config.AccessToken, config.RefreshToken, []byte(config.Cookie))
	if err != nil {
		log.Println(err)
	}
	saveCookie, err := b.SaveCookie()
	config.AccessToken = b.AccessToken
	config.RefreshToken = b.RefreshToken
	config.Username = ""
	config.Password = ""
	config.Cookie = string(saveCookie)
	saveConfig, err := SaveConfig(config)
	if err != nil {
		log.Println(err)
	}
	err = ioutil.WriteFile(configFile, saveConfig, 0777)
	if err != nil {
		log.Println(err)
	}
	log.Println("登录成功")

	// 房间号
	roomId, err := api.GetRoomId(shortId)
	if err != nil {
		fmt.Println(err)
	}

	// 新建分发器
	d := NewDispatcher(roomId)
	d.Add(worker.Danmaku())                  // 显示弹幕
	d.Add(worker.GiftAcknowledge(roomId, b)) // 感谢礼物

	err = d.Connect()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("连接弹幕服务器成功")
	}

	// 开始监听
	go d.Listen()

	reader := bufio.NewReader(os.Stdin)

	// 发送弹幕
	emptyCount := 0
	for emptyCount < 2 {
		text, _ := reader.ReadString('\n')
		content := strings.Trim(text, " \r\n")
		if len(content) > 0 {
			err := b.SendLiveMessage(roomId, content)
			if err != nil {
				log.Println(err)
			} else {
				log.Println("弹幕发送成功", content)
			}
			emptyCount = 0
		} else {
			emptyCount++
		}
	}
}
