package main

import (
	"github.com/gorilla/websocket"
	"gowebsocket.com/impl"
	"net/http"
	"time"
)

//升级为websocket的工具定义
var (
	upgrader = websocket.Upgrader{
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
)

/**
  起一个WS连接
  客户端发什么，服务端就返回什么
*/
func wsHandler(w http.ResponseWriter, r *http.Request) {
	//升级HTTP为WS，完成握手,在response 加入了upgrade:websocket
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//出错直接return
		return
	}

	conn, err := impl.CreateConnection(wsConn)

	//每隔一秒给客户端发送一个心跳消息
	go func() {
		for {
			if beatErr := conn.WriteMessage([]byte("heart beat")); beatErr != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}

	}()

	for {
		data, readErr := conn.ReadMessage()
		if readErr != nil {
			goto ERR
		}
		writeErr := conn.WriteMessage(data)
		if writeErr != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()

}

func main() {

	//set url for visiting: it is: http://127.0.0.1:7777/ws
	http.HandleFunc("/ws", wsHandler)
	//create listen and serve for 0.0.0.0:7000 , do nothing
	_ = http.ListenAndServe("127.0.0.1:7777", nil)

}
