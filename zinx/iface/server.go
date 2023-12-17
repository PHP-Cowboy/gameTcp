package iface

type Server interface {
	//启动
	Start()
	//运行
	Serve()
	//停止
	Stop()
}
