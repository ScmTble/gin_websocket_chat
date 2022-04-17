package main

import (
	"chat/global"
	"chat/hub"
	"chat/log"
	"chat/router"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

func main() {
	log.InitZap()
	global.Upgrade = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// 启动Hub
	hub.NewHub()
	go hub.H.Run()
	e := router.NewRouter()
	fmt.Println(e.Run(":8080"))
}
