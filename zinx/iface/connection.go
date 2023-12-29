package iface

import "net"

// 定义连接接⼝口
type Connection interface {
	Start()
	Stop()
	GetConnId() uint32
	GetTCPConnection() *net.TCPConn
	RemoteAddr() net.Addr
	SendMsg(msgId uint32, data []byte) (err error)
	SendBuffMsg(msgId uint32, data []byte) (err error) //直接将Message数据发送给远程的TCP客户端(有缓冲)
}

type HandleFunc func(conn *net.TCPConn, data []byte, cnt int) (err error)
