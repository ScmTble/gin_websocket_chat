package hub

import (
	"encoding/json"
	"fmt"
)

// Msg 消息格式

type MsgType uint

const (
	TextMsg    MsgType = 1
	NoteMsg    MsgType = 2
	OnlineMsg  MsgType = 3
	OfflineMsg MsgType = 4
)

type Msg struct {
	Type MsgType `json:"type"`
	Msg  string  `json:"msg"`
	From uint    `json:"from"`
	To   uint    `json:"to"`
}

// NewTextMsg 创建一条消息
func NewTextMsg(msg string, to uint, from uint) []byte {
	m := &Msg{
		Type: TextMsg,
		Msg:  msg,
		From: from,
		To:   to,
	}
	marshal, err := json.Marshal(m)
	if err != nil {
		return []byte("json编码错误")
	}
	return marshal
}

// NewNoteMsg 创建一条通知消息
func NewNoteMsg(msg string) []byte {
	m := &Msg{
		Type: NoteMsg,
		Msg:  msg,
		From: 0,
		To:   0,
	}
	marshal, err := json.Marshal(m)
	if err != nil {
		return []byte("json编码错误")
	}
	return marshal
}

// NewOnlineMsg 创建上线消息
func NewOnlineMsg(uid uint) []byte {
	m := &Msg{
		Type: OnlineMsg,
		Msg:  fmt.Sprintf("用户 %d 上线", uid),
		From: 0,
		To:   uid,
	}
	marshal, err := json.Marshal(m)
	if err != nil {
		return []byte("json编码错误")
	}
	return marshal
}
