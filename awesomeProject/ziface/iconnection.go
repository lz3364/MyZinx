package ziface

import "net"

type IConnection interface {

	//启动连接
	Start()

	//停止连接
	Stop()

	//从当前连接获取原始的socket TCPConn
	GetTCPConnection() *net.TCPConn

	//获取当前连接ID
	GetConnID() uint32

	//获取远程客户端地址信息
	RemoteAddr() net.Addr

	SendMsg(msgId uint32, data []byte) error

	SendBuffMsg(msgId uint32, data []byte) error
}

//定义一个统一的业务处理接口
//作为一个回调函数存在
type HandFunc func(*net.TCPConn, []byte, int) error
