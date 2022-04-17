package router

import (
	"chat/api"
	"chat/log"
	"chat/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.GinLogger(log.Logger), middleware.GinRecovery(log.Logger, true))
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "SUCCESS")
	})
	engine.POST("/user/register", api.UserRegister)
	engine.GET("/ws", api.WsHandler)
	return engine
}
