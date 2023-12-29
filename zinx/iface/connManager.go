package iface

type ConnManager interface {
	Add(conn Connection)                   //添加链接
	Remove(conn Connection)                //删除连接
	Get(connId uint32) (Connection, error) //根据ConnID获取链接
	Len() int                              //获取当前连接总数量
	ClearConn()                            //删除并停止所有链接
}
