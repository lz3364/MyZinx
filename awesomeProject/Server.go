package main

import (
	"awesomeProject/ziface"
	"awesomeProject/znet"
	"fmt"
)

type PingRounter struct {
	znet.BaseRounter
}

func (pr *PingRounter) PreHandle(request ziface.IRequest) {
}

func (pr *PingRounter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRounter Handler...")
	fmt.Println("recv from client:msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println("request err:", err)
	}
	return
}

func (pr *PingRounter) PostHandle(request ziface.IRequest) {
}

type HelloZinxRounter struct {
	znet.BaseRounter
}

func (this *HelloZinxRounter) PreHandle(request ziface.IRequest) {
}

func (this *HelloZinxRounter) PostHandle(request ziface.IRequest) {
}

func (this *HelloZinxRounter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRounter handle...")

	fmt.Println("recv from client: msgid = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("Hello Zinx Rounter 0.6"))
	if err != nil {
		fmt.Println(err)
	}
	return
}

func main() {
	//s := znet.NewServer("Zinx[v0.6]")
	s := znet.NewServer()
	s.AddRounter(0, &PingRounter{})
	s.AddRounter(1, &HelloZinxRounter{})
	s.Serve()
}
