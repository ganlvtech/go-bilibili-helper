package websocket

import (
	"bytes"
	"encoding/binary"
)

type BilibiliPacket struct {
	PacketLen uint32
	HeaderLen uint16
	Ver       uint16
	Op        uint32
	Seq       uint32
	Data      []byte
}

func Encode(data []byte, op int) []byte {
	packetLen := 16 + len(data)
	buf := bytes.NewBuffer(make([]byte, 0, packetLen))
	binary.Write(buf, binary.BigEndian, uint32(packetLen))
	binary.Write(buf, binary.BigEndian, uint16(16))
	binary.Write(buf, binary.BigEndian, uint16(2))
	binary.Write(buf, binary.BigEndian, uint32(op))
	binary.Write(buf, binary.BigEndian, uint32(1))
	binary.Write(buf, binary.BigEndian, data)
	return buf.Bytes()
}

func Decode(data []byte) []BilibiliPacket {
	var result []BilibiliPacket
	buf := bytes.NewBuffer(data)
	for i := 0; buf.Len() > 0; i++ {
		packet := BilibiliPacket{}
		binary.Read(buf, binary.BigEndian, &packet.PacketLen)
		binary.Read(buf, binary.BigEndian, &packet.HeaderLen)
		binary.Read(buf, binary.BigEndian, &packet.Ver)
		binary.Read(buf, binary.BigEndian, &packet.Op)
		binary.Read(buf, binary.BigEndian, &packet.Seq)
		dataLen := packet.PacketLen - uint32(packet.HeaderLen)
		packet.Data = make([]byte, dataLen)
		binary.Read(buf, binary.BigEndian, &packet.Data)
		result = append(result, packet)
	}
	return result
}
