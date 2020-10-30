package eqsWebsocket

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	wsConn *websocket.Conn
	inChan chan []byte
	outChan chan []byte
	closeChan chan byte
	mutex sync.Mutex
	isClosed bool
}

// 读取Api
func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <- conn.inChan:
	case <- conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// 发送Api
func (conn *Connection) WriteMessage(data []byte) (err error)  {
	select {
	case conn.outChan <- data:
	case <- conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// 关闭连接的Api
func (conn *Connection) Close()  {
	// 线程安全的Close，可以并发多次调用也叫做可重入的Close
	conn.wsConn.Close()
	conn.mutex.Lock()
	if !conn.isClosed {
		// 关闭chan,但是chan只能关闭一次
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()

}

// 初始化长连接
func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error)  {
	conn = &Connection{
		wsConn: wsConn,
		inChan: make(chan []byte, 1000),
		outChan: make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
	}

	// 启动读协程
	go conn.readLoop()

	// 启动写协程
	go conn.writeLoop()

	return
}

// 内部实现
func (conn *Connection) readLoop()  {
	var (
		data []byte
		err error
	)
	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			goto ERR
		}
		// 容易阻塞到这里，等待inChan有空闲的位置
		select {
		case conn.inChan <- data:
		case <- conn.closeChan: // closeChan关闭的时候执行
			goto ERR
		}
	}

ERR:
	conn.Close()
}

func (conn *Connection) writeLoop()  {
	var (
		data []byte
		err error
	)
	for {
		select {
		case data = <- conn.outChan:
		case <- conn.closeChan:
			goto ERR
		}
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.Close()
}