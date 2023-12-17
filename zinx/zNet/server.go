package zNet

import (
	"fmt"
	"net"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP :%s, Port %d, is starting n", s.IP, s.Port)
	// 监听地址和端口
	serverAddr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr error :", err)
		return
	}

	var listener *net.TCPListener

	// 创建TCP监听器
	listener, err = net.ListenTCP("tcp", serverAddr)
	if err != nil {
		fmt.Println("ListenTCP error:", err)
		return
	}

	defer listener.Close()

	fmt.Println("start ZInx server success")

	var conn *net.TCPConn
	for {
		// 接受客户端连接
		conn, err = listener.AcceptTCP()
		if err != nil {
			fmt.Println("AcceptTCP error:", err)
			continue
		}

		fmt.Println("Accepted new connection from", conn.RemoteAddr())
		go func() {
			// 读取客户端发送的数据
			buffer := make([]byte, 1024)
			var n int
			n, err = conn.Read(buffer)
			if err != nil {
				fmt.Println("Read error:", err)
				return
			}

			fmt.Printf("recv client buf %s, cnt %d\n", buffer, n)

			// 将数据返回给客户端
			_, err = conn.Write(buffer[:n])
			if err != nil {
				fmt.Println("Write error:", err)
				return
			}
		}()
	}

}

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	s.Start()
	//阻塞
	select {}
}

// 初始化Server
func NewServer(name string, version string, ip string, port int) *Server {
	return &Server{
		Name:      name,
		IPVersion: version,
		IP:        ip,
		Port:      port,
	}
}
