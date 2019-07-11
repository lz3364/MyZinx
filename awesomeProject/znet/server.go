package znet

import (
	"awesomeProject/utils"
	"awesomeProject/ziface"
	"errors"
	"fmt"
	"net"
	"time"
)

type Server struct {
	//服务器名称
	Name string

	//服务器网络:TCP4 OR OTHER
	IPVersion string

	//服务器IP
	IP string

	//服务器端口
	Port int

	//当前Server由用户绑定的回调router,也就是Server注册的链接对应的处理
	//业务
	//Rounter ziface.IRounter

	//当前Server的消息管理模块,用来绑定MsgID和对应的处理方法
	msgHandler ziface.IMsgHandle

	//链接管理器
	ConnMgr ziface.IConnManager

	OnConnStart func(connection ziface.IConnection)
	OnConnStop  func(connection ziface.IConnection)
}

func (s *Server) SetOnConnStart(hookfunc func(ziface.IConnection)) {
	s.OnConnStart = hookfunc
}

func (s *Server) SetOnConnStop(hookfunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookfunc
}

func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(connection)
	}
}

func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(connection)
	}
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

/*
 * 回写客户端显示的回调函数
 * 在IConnection中定义的处理函数
 */
func CallbackToClient(conn *net.TCPConn, buf []byte, n int) error {
	fmt.Println("[Conn Handle] Callback to Client...")
	_, err := conn.Write(buf[:n])
	if err != nil {
		fmt.Println("Write buf to client err:", err)
		return errors.New("CallbackToClient error")
	}

	return nil
}

func (s *Server) Start() {
	fmt.Printf("[START] Server[%s] listener at IP:%s Port:%d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)

	// 启动服务后
	go func() {
		s.msgHandler.StartWorkerPool()

		//获取远程TCP地址
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}

		// 监听TCP端口
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen error:", err)
			return
		}

		fmt.Println("start Zinx server...", s.Name, " succ, now listenning...")

		//connection 中的连接ID
		var cid uint32
		cid = 0

		for {

			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept tcp error:", err)
				continue
			}

			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			// 生成新的连接
			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			//ID自增
			cid++

			// 起一个线程，实现对连接的处理
			go dealConn.Start()
		}

	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server, name ", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()
	for {
		time.Sleep(10 * time.Second)
	}
}

func (s *Server) AddRounter(msgId uint32, router ziface.IRounter) {
	s.msgHandler.AddRounter(msgId, router)
	return
}

func NewServer() ziface.IServer {

	//获取配置文件内容
	utils.GlobalObject.Reload()

	s := &Server{
		utils.GlobalObject.Name,
		"tcp4",
		utils.GlobalObject.Host,
		utils.GlobalObject.TcpPort,
		NewMsgHandle(),
		NewConnManager(),
		nil,
		nil,
	}

	return s

}
