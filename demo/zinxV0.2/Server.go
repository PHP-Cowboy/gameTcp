package main

import "gameTcp/zinx/zNet"

func main() {
	//1 创建一个server 句柄 s
	s := zNet.NewServer("[zinx V0.2]", "tcp4", "0.0.0.0", 8090)

	//2 开启服务
	s.Serve()
}
