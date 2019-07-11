package main

import (
	"awesomeProject/utils"
	"awesomeProject/znet"
	"fmt"
	"io"
	"net"
)

func main() {
	utils.GlobalObject.Reload()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", utils.GlobalObject.Host, utils.GlobalObject.TcpPort))
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		go func(conn net.Conn) {
			dp := znet.NewDataPack()
			for {
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println("io reader err:", err)
					break
				}

				message, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("unpack error:", err)
					break
				}

				if message.GetDataLen() > 0 {
					msg := message.(*znet.Message)
					msg.Data = make([]byte, msg.GetDataLen())
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server unpack err:", err)
						return
					}

					fmt.Println("===> recv msg: ID=", msg.Id, " len = ", msg.DataLen, " data=", string(msg.Data))
				}

			}
		}(conn)
	}
}
