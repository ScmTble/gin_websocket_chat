package hub

// Hub 服务器处理中心
type Hub struct {
	// 广播消息通道
	Broadcast chan string

	Clients map[uint]*Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast: make(chan string),
		Clients:   make(map[uint]*Client),
	}
}

// Run Hub运行
func (h *Hub) Run() {
	for {
		select {
		case msg := <-h.Broadcast:
			for _, c := range h.Clients {
				c.Send <- msg
			}
		}
	}
}

// Del 删除客户端
func (h *Hub) Del(uid uint) {
	delete(h.Clients, uid)
}
