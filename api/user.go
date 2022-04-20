package api

import (
	"chat/global"
	"chat/hub"
	"chat/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
		client := hub.NewClient(conn.RemoteAddr().String(), uid, conn, uint64(time.Now().Unix()))

		go client.Write()
		go client.Read()
		// 有用户上线
		hub.H.Register <- client
	}()
}
