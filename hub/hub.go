package hub

import (
	"chat/log"
	"go.uber.org/zap"
)

var H *Hub

// Hub 服务器处理中心
type Hub struct {
	// 广播消息通道
	Broadcast chan *Message
	// 客户端集合
	Clients map[uint]*Client
	// 注册通道
	Register chan *Client
	// 离线通道
	UnRegister chan *Client
}

func NewHub() {
	H = &Hub{
		Broadcast:  make(chan *Message),
		Clients:    make(map[uint]*Client),
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
	}
}

// Run Hub运行
func (h *Hub) Run() {
	for {
		select {
		// 广播消息
		case msg := <-h.Broadcast:
			for _, c := range h.Clients {
				c.Send <- msg
			}
		// 用户连接进来
		case c := <-h.Register:
			//log.Logger.Info("新用户连接", zap.Any("uid", c.Uid))
			h.Clients[c.Uid] = c
		// 用户离开
		case c := <-h.UnRegister:
			h.del(c.Uid)
		}
	}
}

// Del 删除客户端,并关闭连接
func (h *Hub) del(uid uint) {
	_, loaded := h.Clients[uid]
	if loaded {
		delete(h.Clients, uid)
		log.Logger.Debug("断开与客户端连接", zap.Uint("uid", uid))
	}
}
