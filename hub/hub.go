package hub

import (
	"chat/log"
	"sync"
)

var H *Hub

// Hub 服务器处理中心
type Hub struct {
	// 广播消息通道
	Broadcast chan *Message
	// 读写锁
	ClientsLock sync.RWMutex
	// 客户端集合
	Clients map[string]*Client
	// 注册通道
	Register chan *Client
	// 离线通道
	UnRegister chan *Client
}

// NewHub 创建一个新的服务器处理中心
func NewHub() {
	H = &Hub{
		Broadcast:  make(chan *Message, 100),
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client, 100),
		UnRegister: make(chan *Client, 100),
	}
}

// Run Hub运行
func (h *Hub) Run() {
	for {
		select {
		// 用户连接进来
		case c := <-h.Register:
			h.EventRegister(c)
			h.Broadcast <- NewNoteMsg("有人上线", c.UserId)
		// 用户离开
		case c := <-h.UnRegister:
			h.EventUnRegister(c)
			// 广播消息
		case msg := <-h.Broadcast:
			h.ClientsRange(func(uid string, client *Client) (result bool) {
				client.Send <- msg
				return true
			})
		}
	}
}

// EventRegister 用户连接
func (h *Hub) EventRegister(c *Client) {
	log.Logger.Debug("add user " + c.UserId)
	h.AddClient(c)
}

// AddClient 添加用户
func (h *Hub) AddClient(c *Client) {
	h.ClientsLock.Lock()
	defer h.ClientsLock.Unlock()
	h.Clients[c.UserId] = c
}

// EventUnRegister 用户离线
func (h *Hub) EventUnRegister(c *Client) {
	log.Logger.Debug("user leave " + c.UserId)
	h.DelClient(c)
}

// DelClient 删除用户
func (h *Hub) DelClient(c *Client) {
	h.ClientsLock.Lock()
	defer h.ClientsLock.Unlock()

	if _, ok := h.Clients[c.UserId]; ok {
		delete(h.Clients, c.UserId)
	}
}

// ClientsRange 遍历
func (h *Hub) ClientsRange(f func(uid string, client *Client) (result bool)) {

	h.ClientsLock.RLock()
	defer h.ClientsLock.RUnlock()

	for key, value := range h.Clients {
		result := f(key, value)
		if result == false {
			return
		}
	}
	return
}
