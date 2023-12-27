package zNet

import (
	"fmt"
	"gameTcp/zinx/iface"
	"gameTcp/zinx/utils"
	"net"
)

type Server struct {
	Name       string
	IPVersion  string
	Host       string
	Port       int
	MsgHandler iface.MsgHandler
}

// 初始化Server
func NewServer() *Server {
	global := utils.Global
	return &Server{
		Name:       global.Name,
		IPVersion:  global.IPVersion,
		Host:       global.Host,
		Port:       global.TcpPort,
		MsgHandler: NewMsgHandler(),
	}
}

func (s *Server) AddRouter(msgId uint32, router iface.Router) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add router success! msgId = ", msgId)
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at Host :%s, Port %d, is starting n", s.Host, s.Port)

	go func() {
		s.MsgHandler.StartWorkerPool()
		// 监听地址和端口
		serverAddr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.Host, s.Port))
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

		var connId uint32 = 1
		for {

			// 接受客户端连接
			conn, err = listener.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTCP error:", err)
			}

			dealConn := NewConnection(conn, connId, s.MsgHandler)

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
