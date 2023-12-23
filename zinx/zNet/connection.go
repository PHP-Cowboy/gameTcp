package zNet

import (
	"errors"
	"fmt"
	"gameTcp/zinx/iface"
	"io"
	"net"
)

type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnId uint32
	//当前连接的关闭状态
	IsClosed bool
	//该连接的处理方法router
	Router iface.Router
	//告知该链接已经退出/停止的channel
	ExitBuffChan chan struct{}
}

func NewConnection(conn *net.TCPConn, connId uint32, router iface.Router) *Connection {
	return &Connection{
		Conn:         conn,
		ConnId:       connId,
		IsClosed:     false,
		Router:       router,
		ExitBuffChan: make(chan struct{}, 1),
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is  running")
	defer fmt.Println(c.Conn.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		pack := NewPack()

		headData := make([]byte, pack.GetHeadLen())
		//读取我们最大的数据到buf中

		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("receive buf err ", err)
			c.ExitBuffChan <- struct{}{}
			continue
		}

		msg, err := pack.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			c.ExitBuffChan <- struct{}{}
			continue
		}

		var data []byte

		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())

			if _, err = io.ReadFull(c.GetTCPConnection(), data); err != nil {
				return
			}
		}

		msg.SetData(data)

		req := Request{
			Conn: c,
			Msg:  msg,
		}

		go func(req *Request) {
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.AfterHandle(req)
		}(&req)
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

func (c *Connection) SendMsg(msgId uint32, data []byte) (err error) {
	if c.IsClosed {
		return errors.New("Connection closed when send msg")
	}

	pack := NewPack()

	msg, err := pack.Pack(NewMsg(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}

	if _, err = c.Conn.Write(msg); err != nil {
		fmt.Println("Write msg id ", msgId, " error ")
		c.ExitBuffChan <- struct{}{}
		return errors.New("conn Write error")
	}

	return
}
