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

// 创建连接的时候执行
func DoConnectionBegin(conn iface.Connection) {
	fmt.Println("DoConnectionBegin is Called ... ")
	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

// 连接断开的时候执行
func DoConnectionLost(conn iface.Connection) {
	fmt.Println("DoConnectionLost is Called ... ")
}

func main() {
	//1 创建一个server 句柄 s
	s := zNet.NewServer()

	//2 注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//3 添加路由
	s.AddRouter(1, &PingRouter{})
	s.AddRouter(2, &HelloRouter{})

	//4 开启服务
	s.Serve()
}
