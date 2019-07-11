package znet

import "awesomeProject/ziface"

type Request struct {
	conn ziface.IConnection
	//data []byte
	msg ziface.IMessage
}

/*
 * 获取请求的连接
 */
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

/*
 * 获取请求的数据
 */
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
