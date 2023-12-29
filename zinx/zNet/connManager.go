package zNet

import (
	"errors"
	"gameTcp/zinx/iface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]iface.Connection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]iface.Connection),
	}
}

// 添加链接
func (m *ConnManager) Add(conn iface.Connection) {
	m.connLock.Lock()
	defer m.connLock.Unlock()

	_, ok := m.connections[conn.GetConnId()]

	if !ok {
		m.connections[conn.GetConnId()] = conn
	}

	return
}

// 删除连接
func (m *ConnManager) Remove(conn iface.Connection) {
	m.connLock.Lock()
	defer m.connLock.Unlock()

	//删除连接信息
	delete(m.connections, conn.GetConnId())
}

// 根据ConnID获取链接
func (m *ConnManager) Get(connId uint32) (iface.Connection, error) {
	//保护共享资源Map 加读锁
	m.connLock.RLock()
	defer m.connLock.RUnlock()

	conn, ok := m.connections[connId]

	if ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

// 获取当前连接总数量
func (m *ConnManager) Len() int {
	return len(m.connections)
}

// 删除并停止所有链接
func (m *ConnManager) ClearConn() {
	m.connLock.Lock()
	defer m.connLock.Unlock()

	for id, conn := range m.connections {
		conn.Stop()
		delete(m.connections, id)
	}
	return
}
