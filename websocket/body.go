package websocket

import (
	"encoding/binary"

	"github.com/bitly/go-simplejson"
)

func EncodeJoinRoom(roomId int) []byte {
	json2 := simplejson.New()
	json2.Set("roomid", roomId)
	data, _ := json2.Encode()
	return data
}

func DecodeJSON(data []byte) *simplejson.Json {
	json2, _ := simplejson.NewJson(data)
	return json2
}

func DecodePopular(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}
