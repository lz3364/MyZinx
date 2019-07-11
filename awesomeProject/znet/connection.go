package znet

import (
	"awesomeProject/utils"
	"awesomeProject/ziface"
	"errors"
	"fmt"
	"io"
	"net"
)

/*
 * 对每个客户端连接设置属性
 * 包括:
 *   连接的套接字 net.TCPConn
 *   连接的ID
 *   连接是否关闭
 *   连接的处理函数
 *   连接退出的管道
 */
type Connection struct {
	TCPServer ziface.IServer

	//当前连接的socket TCP套接字
	Conn *net.TCPConn

	ConnID uint32

	isClosed bool

	//handleAPI ziface.HandFunc
	// 链接的路由处理方法
	//Rounter ziface.IRounter
	MsgHandle ziface.IMsgHandle

	ExitBuffChan chan bool

	//增加消息缓冲区,用于读写分离
	msgChan chan []byte

	msgBuffChan chan []byte
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack msg error, msg id = ", msgId)
		return errors.New("Pack error msg")
	}

	c.msgBuffChan <- msg

	return nil
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack msg error, msg id = ", msgId)
		return errors.New("Pack error msg")
	}

	c.msgChan <- msg

	return nil
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) (connection *Connection) {

	connection = &Connection{
		server,
		conn,
		connID,
		false,
		msgHandle,
		make(chan bool, 1),
		make(chan []byte),
		make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
	}

	connection.TCPServer.GetConnMgr().Add(connection) //将当前新创建的连接添加到ConnManager中

	return
}

/*
 * 获取远程连接的地址
 */
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

/*
 * 关闭连接，结束当前连接状态
 */
func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}

	c.TCPServer.CallOnConnStop(c)

	c.isClosed = true
	c.Conn.Close()
	c.TCPServer.GetConnMgr().Remove(c)

	c.ExitBuffChan <- true

	close(c.ExitBuffChan)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

/*
 * 处理conn读数据的Goroutine
 */
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println(c.RemoteAddr().String(), "conn reader exit")
	defer c.Stop()

	for {

		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		// IO读满了headData大小的数据后,再读取bodyData的数据
		// client端只发一次,Server端接收两次
		// 第一次读取
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			fmt.Println("Read headdata error:", err)
			c.ExitBuffChan <- true
			continue
		}

		message, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error:", err)
			c.ExitBuffChan <- true
			continue
		}

		var data []byte
		if message.GetDataLen() > 0 {
			data = make([]byte, message.GetDataLen())
			// 第二次读取, 读取的是bodyData的数据
			_, err := io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("read msg data error:", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		message.SetData(data)

		req := Request{
			conn: c,
			msg:  message,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			c.MsgHandle.DoMsgHandler(&req)
		}
	}
}

func (c *Connection) Start() {
	go c.StartReader()

	go c.StartWriter()

	c.TCPServer.CallOnConnStart(c)

	for {
		select {
		case <-c.ExitBuffChan:
			return
		}
	}

}

/*
 *写消息的Goroutine,将用户的数据发给客户端
 */
func (c *Connection) StartWriter() {

	fmt.Println("[Writer Goroutine Is Running...]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("[Writer] write err:", err)
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("[Writer] write err:", err)
					return
				}
			} else {
				fmt.Println("msgBuffChan is Close!")
				break
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}
