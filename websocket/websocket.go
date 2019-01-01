package websocket

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/ganlvtech/go-bilibili-helper/api"
)

type BilibiliWebSocket struct {
	ws     *websocket.Conn
	C      chan BilibiliPacket
	ticker *time.Ticker
}

func New(shortId int) (*BilibiliWebSocket, error) {
	roomId, err := api.GetRoomId(shortId)
	if err != nil {
		return nil, errors.WithMessage(err, "get room id failed")
	}
	return NewBilibiliWebSocket(roomId)
}

func NewBilibiliWebSocket(roomId int) (*BilibiliWebSocket, error) {
	b := new(BilibiliWebSocket)
	b.C = make(chan BilibiliPacket)

	// Create websocket
	var err error
	b.ws, _, err = websocket.DefaultDialer.Dial("wss://broadcastlv.chat.bilibili.com:2245/sub", nil)
	if err != nil {
		return nil, errors.WithMessage(err, "connect websocket failed")
	}

	// receive Message
	done := make(chan bool)
	go b.runReceiveMessage(done)

	// Send Join Room Packet
	err = b.sendJoinRoom(roomId)
	if err != nil {
		return nil, errors.WithMessage(err, "send join room packet failed")
	}
	<-done

	// Run Heartbeat
	b.ticker = time.NewTicker(30 * time.Second)
	go b.runHeartbeat()

	return b, nil
}

func (b *BilibiliWebSocket) runReceiveMessage(done chan bool) error {
	for {
		packets, err := b.receive()
		if err != nil {
			return err
		}
		for i := range packets {
			packet := packets[i]
			if packet.Op == 8 {
				close(done)
			} else {
				b.C <- packet
			}
		}
	}
}

func (b *BilibiliWebSocket) runHeartbeat() error {
	for range b.ticker.C {
		err := b.sendHeartbeat()
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BilibiliWebSocket) send(body []byte, op int) error {
	return b.ws.WriteMessage(websocket.TextMessage, Encode(body, op))
}

func (b *BilibiliWebSocket) sendJoinRoom(roomId int) error {
	return b.send(EncodeJoinRoom(roomId), 7)
}

func (b *BilibiliWebSocket) sendHeartbeat() error {
	return b.send([]byte{}, 2)
}

func (b *BilibiliWebSocket) receive() ([]BilibiliPacket, error) {
	_, message, err := b.ws.ReadMessage()
	if err != nil {
		return nil, err
	}
	packets := Decode(message)
	return packets, nil
}

func (b *BilibiliWebSocket) Close() {
	b.ticker.Stop()
	b.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	b.ws.Close()
}
