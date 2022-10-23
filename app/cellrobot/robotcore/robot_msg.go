/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-20 11:24:15
 * @LastEditTime: 2022-09-20 11:40:11
 * @FilePath: \trcell\app\cellrobot\robotcore\robot_msg.go
 */
package robotcore

import (
	"encoding/binary"
	"errors"
)

//const ClientHeaderSize = 12
const ClientHeaderSize = 6

type ClientMsgHeader struct {
	Len      uint32
	MainType uint8
	SubType  uint8
	ID       uint32 // 临时值，不参与序列化
	Result   uint16 // 临时值，不参与序列化
}

func NewClientHeader(mainType uint8, subType uint8, id uint32) *ClientMsgHeader {
	return &ClientMsgHeader{
		MainType: mainType,
		SubType:  subType,
		ID:       id,
	}
}

func (head *ClientMsgHeader) Encode() []byte {
	buf := make([]byte, ClientHeaderSize)
	binary.LittleEndian.PutUint32(buf, head.Len)
	buf[4] = head.MainType
	buf[5] = head.SubType
	// binary.LittleEndian.PutUint32(buf[6:], head.ID)
	// binary.LittleEndian.PutUint16(buf[10:], head.Result)
	return buf
}

func (head *ClientMsgHeader) Decode(data []byte) {
	head.Len = binary.LittleEndian.Uint32(data[0:4])
	head.MainType = uint8(data[4])
	head.SubType = uint8(data[5])
	// head.ID = binary.LittleEndian.Uint32(data[6:])
	// head.Result = binary.LittleEndian.Uint16(data[10:])
}

type ClientMessage struct {
	Head *ClientMsgHeader
	Data []byte
}

func NewClientMessage(header *ClientMsgHeader) *ClientMessage {
	msg := &ClientMessage{
		Head: header,
		Data: []byte{},
	}
	return msg
}

func (msg *ClientMessage) Encode() []byte {
	msg.Head.Len = uint32(ClientHeaderSize + len(msg.Data))
	hd := msg.Head.Encode()
	return append(hd, msg.Data...)
}

func (msg *ClientMessage) Decode(data []byte) error {
	if len(data) < ClientHeaderSize {
		return errors.New("decode error, len(data) < header size")
	}
	msg.Head.Decode(data[:ClientHeaderSize])
	msg.Data = data[ClientHeaderSize:]
	return nil
}

// 生成一个消息
func MakeClientTcpMessage(mainType uint8, subType uint8, id uint32, result uint16) *ClientMessage {
	hd := NewClientHeader(mainType, subType, id)
	hd.Result = result
	msg := NewClientMessage(hd)
	return msg
}
