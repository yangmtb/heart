package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if nil != err {
		fmt.Printf("failed to upgrade:%+v", err)
		return
	}
	for {
		t, msg, err := conn.ReadMessage()
		if nil != err {
			break
		}
		conn.WriteMessage(t, msg)
	}
}

func readMsgFromRds(uid int64, conn *websocket.Conn) (res *string) {
	t, msg, err := conn.ReadMessage()
	if nil != err {
		return
	}
	fmt.Println(t, msg)
	return
}

func sendMsg2Rds(uid int64, msg string) {
	fmt.Println(uid, msg)
}

func revAndSend(uid int64, conn *websocket.Conn) {
	for {
		if msg := readMsgFromRds(uid, conn); nil != msg {
			conn.WriteMessage(websocket.TextMessage, []byte("pong"))
		} else {
			break
		}
	}
}
