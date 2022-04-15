// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"chat/hub"
	"chat/log"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

func Handle(h *hub.Hub, conn *websocket.Conn, uid uint) {
	c := &hub.Client{
		H:    h,
		Conn: conn,
		Uid:  uid,
		Name: "ScmTble",
		Send: make(chan *hub.Message),
	}
	// 有用户上线时发送广播消息
	h.Broadcast <- hub.NewOnlineMsg(c.Uid)
	h.Clients.Store(c.Uid, c)
	go c.ListenMsg()
	go c.RevMsg()
}

func main() {
	log.InitZap()
	upgrade := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// 启动Hub
	h := hub.NewHub()
	go h.Run()
	// 升级为websocket
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			return
		}
		uid := r.Form["uid"][0]
		parseUint, err := strconv.ParseUint(uid, 10, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		conn, err := upgrade.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		go Handle(h, conn, uint(parseUint))
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		return
	}
}
