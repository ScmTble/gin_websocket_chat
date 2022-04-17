package hub

import (
	"chat/log"
	"errors"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

// Client 每一个连接进来的都是一个Client对象
type Client struct {
	Conn *websocket.Conn
	Uid  uint
	Name string
	// 发送给本User的channel
	Send chan *Message
}

// Write 发送消息给客户端
func (c *Client) Write() {
	// 心跳检查
	ticker := time.NewTicker(time.Second * 3)
	defer func() {
		ticker.Stop()
		H.UnRegister <- c
		c.Conn.Close()
		//log.Logger.Debug("消息处理程序停止", zap.Uint("uid", c.Uid))
	}()
	for {
		select {
		case msg := <-c.Send:
			// 当有消息到达时
			log.Logger.Debug("用户管道收到消息", zap.Uint("uid", c.Uid), zap.Any("msg", msg))
			if err := c.Conn.WriteJSON(msg); err != nil {
				return
			}
		case <-ticker.C:
			// 心跳检查(3s)
			err := c.Conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				//log.Logger.Error("心跳检查错误", zap.Uint("uid", c.Uid), zap.Error(err))
				return
			} else {
				log.Logger.Debug("心跳检查正常", zap.Uint("uid", c.Uid))
			}
		}
	}
}

// Read 处理客户端conn发送过来的消息
func (c *Client) Read() {
	defer func() {
		H.UnRegister <- c
		c.Conn.Close()
		//log.Logger.Debug("消息监听程序停止", zap.Uint("uid", c.Uid))
	}()
	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			//log.Logger.Error("读取客户端conn消息异常", zap.Uint("uid", c.Uid), zap.Error(err))
			break
		}
		msg.From = c.Uid
		log.Logger.Debug("接收到客户端conn发来的消息", zap.Uint("uid", c.Uid), zap.Any("msg", msg))
		if err = c.senMsg(&msg); err != nil {
			log.Logger.Error("发送消息异常", zap.Error(err))
			if err := c.Conn.WriteJSON(NewNoteMsg("消息发送失败!")); err != nil {
				return
			}
		} else {
			if err := c.Conn.WriteJSON(NewNoteMsg("消息发送成功!")); err != nil {
				return
			}
		}
	}
}

// 发送消息
func (c *Client) senMsg(msg *Message) error {
	client, ok := H.Clients[msg.To]
	if !ok {
		return errors.New("找不到对应的客户端")
	}
	client.Send <- NewTextMsg(msg.Msg, msg.To, msg.From)
	return nil
}
