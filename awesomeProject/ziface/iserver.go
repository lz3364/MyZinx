package ziface

type IServer interface {
	Start()

	Stop()

	Serve()

	/*
	 * 路由功能,给当前的服务注册一个路由业务方法,
	 * 供客户端处理链接时使用
	 * 提供msgid,可以绑定到路由上
	 */
	AddRounter(msgId uint32, router IRounter)

	GetConnMgr() IConnManager

	SetOnConnStart(func(IConnection))
	SetOnConnStop(func(connection IConnection))
	CallOnConnStart(connection IConnection)
	CallOnConnStop(connection IConnection)
}
