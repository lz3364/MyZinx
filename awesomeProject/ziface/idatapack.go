package ziface

/*
 * 封包数据和拆包数据
 * 直接面向TCP连接中的数据流,为传输数据添加头部信息,用于处理TCP的粘包问题
 * 每条包数据,分为head 和 body两种
 * head中定义了包头信息,以区分和其他包的不同
 * 包头信息包括
 * 	数据包长度
 * 	数据包ID
 * todo:理解什么叫TCP粘包
 */
type IDataPack interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	Unpack([]byte) (IMessage, error)
}
