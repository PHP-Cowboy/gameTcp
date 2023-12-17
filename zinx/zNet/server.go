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

func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) (err error) {
	fmt.Println("[Conn Handle] CallBackToClient ... ")
	_, err = conn.Write(data[:cnt])
	if err != nil {
		fmt.Println("Write error:", err)
		return
	}
	return
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP :%s, Port %d, is starting n", s.IP, s.Port)

	go func() {
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

		var connId uint32 = 0
		for {

			// 接受客户端连接
			conn, err = listener.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTCP error:", err)
			}

			dealConn := NewConnection(conn, connId, CallBackToClient)

			connId++

			fmt.Println("Accepted new connection from", conn.RemoteAddr())

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] zInx server , name ", s.Name)
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
