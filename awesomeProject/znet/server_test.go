package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func ClientTest()  {
	fmt.Println("client test...start")
	time.Sleep(10*time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil{
		fmt.Println("dial err:", err)
		return
	}
	_, err = conn.Write([]byte("hello world I am zinx"))
	if err != nil{
		fmt.Println("write error:", err)
		return
	}
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil{
		fmt.Println("conn read err:", err)
		return
	}
	fmt.Println("recv buf:", string(buf[:n]))

	time.Sleep(time.Second)
}

/*
 * 服务器测试
 */
func TestServer(t *testing.T)  {

	s := NewServer("[zinx v0.1]")
	go ClientTest()
	s.Serve()
}
