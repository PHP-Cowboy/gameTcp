package iface

import "net"

// 定义连接接⼝口
type Connection interface {
	Start()
	Stop()
	GetConnId() uint32
}

type HandleFunc func(conn *net.TCPConn, data []byte, cnt int) (err error)
