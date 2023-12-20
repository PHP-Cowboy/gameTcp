package zNet

import "gameTcp/zinx/iface"

type Request struct {
	Conn iface.Connection
	Data []byte
}

func (r *Request) GetConnection() iface.Connection {
	return r.Conn
}

func (r *Request) GetData() []byte {
	return r.Data
}
