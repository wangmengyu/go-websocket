package impl

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

//建立连接对象：长连接本身，读取的通道，写出的通道
type Connection struct {
	wsConn    *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte
	isClose   bool
	mutex     sync.Mutex
}

/**创建连接*/
func CreateConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:    wsConn,
		inChan:    make(chan []byte, 10),
		outChan:   make(chan []byte, 10),
		closeChan: make(chan byte, 1),
		isClose:   false,
	}
	//开启读写协程
	go conn.readLoop()
	go conn.writeLoop()
	return
}

/**
从inChan中读取数据并返回
*/
func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("conn is closed")
	}
	return
}

/**
写入数据到outChan
*/
func (conn *Connection) WriteMessage(data []byte) error {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		return errors.New("conn is closed")
	}
	return nil
}

/**
关闭连接，线程安全，可重入的Close
*/
func (conn *Connection) Close() {
	conn.wsConn.Close()

	//关闭掉closeChan,确保线程安全，需要加入mutex锁
	conn.mutex.Lock()
	if conn.isClose == false {
		close(conn.closeChan)
		conn.isClose = true
	}
	conn.mutex.Unlock()
}

/**
  不断的把读取到的数据，放入inChan
*/
func (conn *Connection) readLoop() {
	for {
		_, data, err := conn.wsConn.ReadMessage()
		if err != nil {
			conn.Close()
			break
		}
		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			//阻塞时候回走到这里，执行关闭连接
			conn.Close()
			break
		}

	}
}

/**
  不断从outChan取得数据，发送出去
*/

func (conn *Connection) writeLoop() {
	for {
		select {
		case data := <-conn.outChan:
			err := conn.wsConn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				conn.Close()
				break
			}
		case <-conn.closeChan:
			conn.Close()
			break
		}

	}
}
