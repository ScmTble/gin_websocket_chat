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

// Message 消息格式
type Message struct {
	MsgType MsgType `json:"type"`    //消息类型
	Message string  `json:"message"` //要发送的内容
	UserID  string  `json:"userID"`  //对方id
}

// NewTextMsg 创建一条消息
func NewTextMsg(msg string, userID string) *Message {
	m := &Message{
		MsgType: TextMsg,
		Message: msg,
		UserID:  userID,
	}
	return m
}

// NewNoteMsg 创建一条通知消息
func NewNoteMsg(msg string, uid string) *Message {
	m := &Message{
		MsgType: NoteMsg,
		Message: msg,
		UserID:  uid,
	}
	return m
}

// NewOnlineMsg 创建上线消息
func NewOnlineMsg(uid uint) *Message {
	m := &Message{
		MsgType: OnlineMsg,
		Message: fmt.Sprintf("用户 %d 上线", uid),
		UserID:  "",
	}
	return m
}
