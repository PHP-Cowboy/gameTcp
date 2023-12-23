package zNet

import "gameTcp/zinx/iface"

type Request struct {
	Conn iface.Connection
	Msg  iface.Message
}

func (r *Request) GetConnection() iface.Connection {
	return r.Conn
}

func (r *Request) GetData() []byte {
	return r.Msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.Msg.GetMsgId()
}
