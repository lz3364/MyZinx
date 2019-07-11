package main

import (
	"awesomeProject/znet"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("Client is starting...")
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial error:", err)
		return
	}

	for {
		/*
			_, err := conn.Write([]byte("Zinx v0.3"))
			if err != nil{
				fmt.Println("client Write error:", err)
				return
			}

			buf := make([]byte, 512)
			n, err := conn.Read(buf)
			if err != nil{
				fmt.Println("client read err:", err)
				return
			}
			fmt.Println("server call back:", string(buf[:n]))

		*/
		dp := znet.NewDataPack()
		msg, err := dp.Pack(znet.NewMsgPackage(1, []byte("Zinx V0.5")))
		if err != nil {
			fmt.Println("Pack Data error:", err)
			return
		}
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("Conn Write error:", err)
			return
		}

		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("read headdata error:", err)
			return
		}

		fmt.Println("Unpack data len: ", len(headData))

		message, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("Unpack error:", err)
			return
		}

		if message.GetDataLen() > 0 {

			//todo:接口的强制转换为实现接口的类,使用.()操作
			//由于message是一个interface,不能直接操作,因此只能转换成结构体
			msg := message.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())
			_, err = io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("read body data error:", err)
				return
			}

			fmt.Println("===> Recv Msg:ID = ", msg.GetMsgId(), " data = ", string(msg.GetData()))
		}

		time.Sleep(time.Second)
	}
}
