package hub

import (
	"chat/log"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"runtime/debug"
)

// Client 每一个连接进来的都是一个Client对象
// 用户连接
type Client struct {
	Addr          string          // 客户端地址
	Socket        *websocket.Conn // 用户连接
	Send          chan *Message   // 待发送的数据
	AppId         uint32          // 登录的平台Id app/web/ios
	UserId        string          // 用户Id，用户登录以后才有
	FirstTime     uint64          // 首次连接事件
	HeartbeatTime uint64          // 用户上次心跳时间
	LoginTime     uint64          // 登录时间 登录以后才有
}

// NewClient 初始化
func NewClient(addr string, uid string, socket *websocket.Conn, firstTime uint64) (client *Client) {
	client = &Client{
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan *Message, 100),
		UserId:        uid,
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
	}

	return
}

// Write 发送消息给客户端
func (c *Client) Write() {
	// panic捕获
	defer func() {
		if a := recover(); a != nil {
			log.Logger.Error("write stop", zap.String("msg", string(debug.Stack())))
		}
	}()
	defer func() {
		H.UnRegister <- c
		c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// c.Send 数据为空，并关闭了
				log.Logger.Info("client发送数据 关闭连接", zap.String("uid", c.UserId))
				return
			}
			c.Socket.WriteJSON(message)
		}
	}
}

// read 处理客户端conn发送过来的消息
func (c *Client) Read() {
	// panic捕获
	defer func() {
		if a := recover(); a != nil {
			log.Logger.Error("read stop", zap.String("msg", string(debug.Stack())))
		}
	}()
	defer func() {
		// 关闭接收通道
		close(c.Send)
	}()
	for {
		var msg Message
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			log.Logger.Error("读取客户端conn消息异常", zap.String("uid", c.UserId), zap.Error(err))
			return
		}
		log.Logger.Debug("接收到客户端conn发来的消息", zap.String("uid", c.UserId), zap.Any("msg", msg))
		ProcessData(c, &msg)
	}
}

// ProcessData 消息处理
func ProcessData(c *Client, msg *Message) {
	return
}
