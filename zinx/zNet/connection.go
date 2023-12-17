package zNet

import (
	"fmt"
	"gameTcp/zinx/iface"
	"net"
)

type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnId uint32
	//当前连接的关闭状态
	IsClosed bool
	//该连接的处理方法api
	HandleApi iface.HandleFunc
	//告知该链接已经退出/停止的channel
	ExitBuffChan chan struct{}
}

func NewConnection(conn *net.TCPConn, connId uint32, handleApi iface.HandleFunc) *Connection {
	return &Connection{
		Conn:         conn,
		ConnId:       connId,
		IsClosed:     false,
		HandleApi:    handleApi,
		ExitBuffChan: make(chan struct{}, 1),
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is  running")
	defer fmt.Println(c.Conn.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		//读取我们最大的数据到buf中
		var buf = make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("receive buf err ", err)
			c.ExitBuffChan <- struct{}{}
			continue
		}
		//调用当前链接业务(这里执行的是当前conn的绑定的handle方法)
		if err = c.HandleApi(c.Conn, buf, cnt); err != nil {
			fmt.Println("connID ", c.ConnId, " handle is error")
			c.ExitBuffChan <- struct{}{}
			return
		}
	}
}

func (c *Connection) Start() {
	go c.StartReader()

	select {
	case <-c.ExitBuffChan:
		//得到退出消息，不再阻塞
		return
	}
}

func (c *Connection) Stop() {
	if c.IsClosed {
		return
	}
	c.IsClosed = true

	err := c.Conn.Close()
	if err != nil {
		return
	}

	//通知从缓冲队列读数据的业务，该链接已经关闭
	c.ExitBuffChan <- struct{}{}

	//关闭该链接全部管道
	close(c.ExitBuffChan)
}

// 获取当前连接ID
func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

// 从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
