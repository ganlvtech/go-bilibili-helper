package main

import (
	"github.com/ganlvtech/go-bilibili-helper/websocket"
)

type Dispatcher struct {
	roomId            int
	bilibiliWebSocket *websocket.BilibiliWebSocket
	Channels          []chan websocket.BilibiliPacket
}

func NewDispatcher(roomId int) *Dispatcher {
	d := new(Dispatcher)
	d.roomId = roomId
	return d
}

func (d *Dispatcher) Connect() error {
	var err error
	d.bilibiliWebSocket, err = websocket.New(d.roomId)
	return err
}

func (d *Dispatcher) Add(c chan websocket.BilibiliPacket) {
	d.Channels = append(d.Channels, c)
}

func (d *Dispatcher) Listen() {
	for packet := range d.bilibiliWebSocket.C {
		for _, c := range d.Channels {
			select {
			case c <- packet:
			default:
			}
		}
	}
}
