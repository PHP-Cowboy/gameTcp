package main

import (
	"fmt"
	"gameTcp/zinx/zNet"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("tcp dail err:", err.Error())
		return
	}

	pack := zNet.NewPack()

	msg1 := zNet.NewMsg(1, []byte{'h', 'e', 'l', 'l', 'o'})

	sendData1, _ := pack.Pack(msg1)

	msg2 := zNet.NewMsg(2, []byte{'w', 'o', 'r', 'l', 'd', '!', '!'})

	sendData2, _ := pack.Pack(msg2)

	sendData1 = append(sendData1, sendData2...)

	//向服务器端写数据
	_, err = conn.Write(sendData1)
	if err != nil {
		return
	}

	//客户端阻塞
	select {}
}
