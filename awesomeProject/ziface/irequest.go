package ziface

/*
 * 将客户端的连接和数据包装到Request里面去
 */
type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
	GetMsgID() uint32
}
