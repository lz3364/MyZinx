package ziface

type IConnManager interface {
	Add(conn IConnection)                   //增加链接
	Remove(conn IConnection)                //移除链接
	Get(connID uint32) (IConnection, error) //通过连接ID获取链接
	Len() int                               //链接数量
	ClearConn()                             //删除并停止所有链接
}
