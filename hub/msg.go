package hub

import (
	"fmt"
)

// Message 消息格式

type MsgType uint

const (
	TextMsg   MsgType = 1
	NoteMsg   MsgType = 2
	OnlineMsg MsgType = 3
)

type Message struct {
	Type MsgType `json:"type"`
	Msg  string  `json:"msg"`
	From uint    `json:"from"`
	To   uint    `json:"to"`
}

// NewTextMsg 创建一条消息
func NewTextMsg(msg string, to uint, from uint) *Message {
	m := &Message{
		Type: TextMsg,
		Msg:  msg,
		From: from,
		To:   to,
	}
	return m
}

// NewNoteMsg 创建一条通知消息
func NewNoteMsg(msg string) *Message {
	m := &Message{
		Type: NoteMsg,
		Msg:  msg,
		From: 0,
		To:   0,
	}
	return m
}

// NewOnlineMsg 创建上线消息
func NewOnlineMsg(uid uint) *Message {
	m := &Message{
		Type: OnlineMsg,
		Msg:  fmt.Sprintf("用户 %d 上线", uid),
		From: 0,
		To:   uid,
	}
	return m
}
