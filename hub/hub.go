package hub

import (
	"chat/log"
	"go.uber.org/zap"
	"sync"
)

// Hub 服务器处理中心
type Hub struct {
	// 广播消息通道
	Broadcast chan *Message
	// 客户端集合
	Clients sync.Map
}

func NewHub() *Hub {
	return &Hub{
		Broadcast: make(chan *Message),
	}
}

// Run Hub运行
func (h *Hub) Run() {
	for {
		select {
		case msg := <-h.Broadcast:
			h.Clients.Range(func(key, value any) bool {
				client := value.(*Client)
				client.Send <- msg
				return true
			})
		}
	}
}

// Del 删除客户端,并关闭连接
func (h *Hub) Del(uid uint) {
	value, loaded := h.Clients.LoadAndDelete(uid)
	if loaded {
		client := value.(*Client)
		if err := client.Conn.Close(); err == nil {
			log.Logger.Debug("断开与客户端连接", zap.Uint("uid", client.Uid))
		} else {
			log.Logger.Debug("关闭连接出错", zap.Error(err))
		}
	}
}
