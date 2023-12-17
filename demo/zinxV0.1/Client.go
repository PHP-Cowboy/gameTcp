package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8090")
	if err != nil {
		fmt.Println("tcp dail err:", err.Error())
		return
	}

	for {
		var cnt int
		_, err = conn.Write([]byte("Hello zInx V0.1"))
		if err != nil {
			fmt.Println("write err:", err.Error())
			return
		}

		buf := make([]byte, 1024)

		cnt, err = conn.Read(buf)
		if err != nil {
			fmt.Println("read err:", err.Error())
			return
		}

		fmt.Printf("server call back: %s, cnt = %d n", buf, cnt)

		time.Sleep(time.Second * 1)
	}
}
