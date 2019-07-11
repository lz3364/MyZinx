package main

import (
	"awesomeProject/utils"
	"awesomeProject/znet"
	"fmt"
	"net"
)

func main() {
	utils.GlobalObject.Reload()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", utils.GlobalObject.Host, utils.GlobalObject.TcpPort))
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}

	defer conn.Close()

	dp := znet.NewDataPack()

	msg1 := &znet.Message{
		Id:      0,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("pack error:", err)
		return
	}

	_, err = conn.Write(sendData1)
	if err != nil {
		fmt.Println("write data error:", err)
		return
	}

	select {}

}
