package main

import (
	"fmt"
	"gameTcp/zinx/iface"
	"gameTcp/zinx/zNet"
)

type PingRouter struct {
	zNet.BaseRouter
}

func (pr *PingRouter) PreHandle(req iface.Request) {
	fmt.Println("Call Router PreHandle")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("before ping ....\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func (pr *PingRouter) Handle(req iface.Request) {
	fmt.Println("Call PingRouter Handle")
	if err := req.GetConnection().SendMsg(1, []byte("ping...ping...ping")); err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func (pr *PingRouter) AfterHandle(req iface.Request) {
	fmt.Println("Call Router PostHandle")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("After ping .....\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func main() {
	//1 创建一个server 句柄 s
	s := zNet.NewServer()

	//2 添加路由
	s.AddRouter(&PingRouter{})

	//3 开启服务
	s.Serve()
}
