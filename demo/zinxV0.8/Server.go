package main

import (
	"fmt"
	"gameTcp/zinx/iface"
	"gameTcp/zinx/zNet"
)

type PingRouter struct {
	zNet.BaseRouter
}

func (r *PingRouter) Handle(req iface.Request) {
	fmt.Println("Call HelloRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("receive from client : msgId=", req.GetMsgId(), ", data=", string(req.GetData()))

	err := req.GetConnection().SendMsg(1, []byte("Hello ZInx Router V0.6"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	zNet.BaseRouter
}

func (r *HelloRouter) Handle(req iface.Request) {
	fmt.Println("Call HelloRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("receive from client : msgId=", req.GetMsgId(), ", data=", string(req.GetData()))

	err := req.GetConnection().SendMsg(2, []byte("Hello ZInx Router V0.6"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//1 创建一个server 句柄 s
	s := zNet.NewServer()

	//2 添加路由
	s.AddRouter(1, &PingRouter{})
	s.AddRouter(2, &HelloRouter{})

	//3 开启服务
	s.Serve()
}
