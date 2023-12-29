package zNet

import (
	"errors"
	"fmt"
	"gameTcp/zinx/iface"
	"gameTcp/zinx/utils"
	"io"
	"net"
	"sync"
)

type Connection struct {
	TcpServer    iface.Server           //当前conn属于哪个server
	Conn         *net.TCPConn           //当前连接的socket TCP套接字
	ConnId       uint32                 //当前连接的ID 也可以称作为SessionID，ID全局唯一
	IsClosed     bool                   //当前连接的关闭状态
	MsgHandler   iface.MsgHandler       //消息管理MsgId和对应处理方法的消息管理模块
	ExitBuffChan chan struct{}          //告知该链接已经退出/停止的channel
	msgChan      chan []byte            //无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan  chan []byte            //有关冲管道，用于读、写两个goroutine之间的消息通信
	property     map[string]interface{} //链接属性
	propertyLock sync.RWMutex           //保护链接属性修改的锁
}

func NewConnection(server iface.Server, conn *net.TCPConn, connId uint32, msgHandler iface.MsgHandler) *Connection {
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnId:       connId,
		IsClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan struct{}, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.Global.MaxMsgChanLen),
		property:     make(map[string]interface{}),
	}

	//将新创建的Conn添加到链接管理中
	c.TcpServer.GetConnManager().Add(c)

	return c
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Data error:, ", err, " Conn Writer exit")
					return
				}
			}
			return
		case <-c.ExitBuffChan:
			//conn已经关闭
			return
		}
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

		if utils.Global.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.Do(&req)
		}

	}
}

func (c *Connection) Start() {
	go c.StartReader()
	go c.StartWriter()
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)

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

	//如果用户注册了该链接的关闭回调业务,调用
	c.TcpServer.CallOnConnStop(c)

	err := c.Conn.Close()
	if err != nil {
		return
	}

	//通知从缓冲队列读数据的业务，该链接已经关闭
	c.ExitBuffChan <- struct{}{}

	//将链接从连接管理器中删除
	c.TcpServer.GetConnManager().Remove(c)

	//关闭该链接全部管道
	close(c.ExitBuffChan)
	close(c.msgBuffChan)
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

	//写回客户端
	c.msgChan <- msg //将之前直接回写给conn.Write的方法 改为 发送给Channel 供Writer读取

	return
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) (err error) {
	if c.IsClosed {
		return errors.New("Connection closed when send msg")
	}

	pack := NewPack()

	msg, err := pack.Pack(NewMsg(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	c.msgBuffChan <- msg //将之前直接回写给conn.Write的方法 改为 发送给Channel 供Writer读取

	return
}

// 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	value, ok := c.property[key]

	if !ok {
		return nil, errors.New("not found")
	}
	return value, nil
}

// 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
