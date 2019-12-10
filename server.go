package main

import (
	"github.com/gorilla/websocket"
	"net/http"
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
	//升级HTTP为WS，完成第一次应答,在response 加入了upgrade:websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//出错直接return
		return
	}
	for {
		//不断读取消息
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			continue
		}
		//不断写入消息
		if err = conn.WriteMessage(msgType, data); err != nil {
			conn.Close()
			continue
		}

	}

}

func main() {

	//set url for visiting: it is: http://127.0.0.1:7777/ws
	http.HandleFunc("/ws", wsHandler)
	//create listen and serve for 0.0.0.0:7000 , do nothing
	_ = http.ListenAndServe("127.0.0.1:7777", nil)

}
