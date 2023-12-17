package main

import "gameTcp/zinx/zNet"

func main() {
	server := zNet.NewServer("zInx0.1", "tcp4", "0.0.0.0", 8090)

	server.Serve()
}
