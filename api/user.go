package api

import (
	"chat/global"
	"chat/hub"
	"chat/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRegister(c *gin.Context) {
	c.JSON(http.StatusOK, "未完成")
}

func WsHandler(c *gin.Context) {
	token := c.Query("token")
	uid, v := utils.VerifyToken(token)
	if !v {
		c.JSON(http.StatusBadRequest, "token验证失败")
		return
	}
	conn, err := global.Upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	go func() {
		c := &hub.Client{
			Conn: conn,
			Uid:  uid,
			Name: "ScmTble",
			Send: make(chan *hub.Message),
		}
		hub.H.Register <- c
		// 有用户上线时发送广播消息
		hub.H.Broadcast <- hub.NewOnlineMsg(c.Uid)
		go c.Read()
		go c.Write()
	}()
}
