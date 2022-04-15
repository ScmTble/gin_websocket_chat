package hub

import (
	"chat/log"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

// Client 每一个连接进来的都是一个Client对象
type Client struct {
	// Hub
	H    *Hub
	Conn *websocket.Conn
	Uid  uint
	Name string

	// 发送给本User的channel
	Send chan string
}

// ListenMsg 监听改用户管道中是否有消息需要发送
func (c *Client) ListenMsg() {
	// 心跳检查
	ticker := time.NewTicker(time.Second * 3)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
		c.H.Del(c.Uid)
		log.Logger.Debug("消息处理程序停止", zap.Uint("uid", c.Uid))
	}()
	for {
		select {
		case msg := <-c.Send:
			// 当有消息到达时
			log.Logger.Debug("用户管道收到消息", zap.Uint("uid", c.Uid), zap.String("msg", msg))
			if err := c.Conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				return
			}
		case <-ticker.C:
			// 心跳检查(3s)
			err := c.Conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Logger.Error("心跳检查错误", zap.Error(err))
				return
			} else {
				log.Logger.Debug("心跳检查正常", zap.Uint("uid", c.Uid), zap.Int("在线人数", len(c.H.Clients)))
			}
		}
	}
}

// RevMsg 处理客户端conn发送过来的消息
func (c *Client) RevMsg() {
	defer func() {
		c.Conn.Close()
		c.H.Del(c.Uid)
		log.Logger.Debug("消息监听程序停止", zap.Uint("uid", c.Uid))
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Logger.Error("读取客户端conn消息异常", zap.Error(err))
			break
		}
		var msg Msg
		// 解码消息
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Logger.Error("Json编码异常", zap.Error(err))
			if err := c.Conn.WriteMessage(websocket.TextMessage, NewNoteMsg("消息格式不对!")); err != nil {
				return
			}
			continue
		}
		msg.From = c.Uid
		log.Logger.Debug("接收到客户端conn发来的消息", zap.Uint("uid", c.Uid), zap.Any("msg", msg))
		if err = c.senMsg(&msg); err != nil {
			log.Logger.Error("发送消息异常", zap.Error(err))
			if err := c.Conn.WriteMessage(websocket.TextMessage, NewNoteMsg("消息发送失败!")); err != nil {
				return
			}
		} else {
			if err := c.Conn.WriteMessage(websocket.TextMessage, NewNoteMsg("消息发送成功!")); err != nil {
				return
			}
		}
	}
}

// 发送消息
func (c *Client) senMsg(msg *Msg) error {
	client, ok := c.H.Clients[msg.To]
	if !ok {
		return errors.New("找不到对应的客户端")
	}
	client.Send <- string(NewTextMsg(msg.Msg, msg.To, msg.From))
	return nil
}
