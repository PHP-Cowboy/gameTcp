package zNet

import (
	"fmt"
	"gameTcp/zinx/iface"
	"gameTcp/zinx/utils"
	"net"
)

type Server struct {
	Name        string                      //服务器的名称
	IPVersion   string                      //tcp4 or other
	Host        string                      //服务绑定的Host
	Port        int                         //服务绑定的端口
	MsgHandler  iface.MsgHandler            //当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	ConnManager iface.ConnManager           //当前Server的链接管理器
	OnConnStart func(conn iface.Connection) //该Server的连接创建时Hook函数
	OnConnStop  func(conn iface.Connection) //该Server的连接断开时的Hook函数
}

// 初始化Server
func NewServer() *Server {
	global := utils.Global
	return &Server{
		Name:        global.Name,
		IPVersion:   global.IPVersion,
		Host:        global.Host,
		Port:        global.TcpPort,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
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
			if s.ConnManager.Len() >= utils.Global.MaxConn {
				continue
			}

			// 接受客户端连接
			conn, err = listener.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTCP error:", err)
			}

			dealConn := NewConnection(s, conn, connId, s.MsgHandler)

			connId++

			fmt.Println("Accepted new connection from", conn.RemoteAddr())

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] zInx server , name ", s.Name)
	s.ConnManager.ClearConn()
}

func (s *Server) Serve() {
	s.Start()
	//阻塞
	select {}
}

func (s *Server) GetConnManager() iface.ConnManager {
	return s.ConnManager
}

// 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(iface.Connection)) {
	s.OnConnStart = hookFunc
}

// 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(iface.Connection)) {
	s.OnConnStop = hookFunc
}

// 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn iface.Connection) {
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}

// 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn iface.Connection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}
