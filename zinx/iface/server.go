package iface

type Server interface {
	//启动服务器方法
	Start()
	//开启业务服务方法
	Serve()
	//停止服务器方法
	Stop()
	//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	AddRouter(msgId uint32, router Router)
	//得到链接管理
	GetConnManager() ConnManager
	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func(Connection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(Connection))
	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn Connection)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn Connection)
}
