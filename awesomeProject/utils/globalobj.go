package utils

import (
	"awesomeProject/ziface"
	"encoding/json"
	"io/ioutil"
)

/*
 * 全局变量
 * 存储一切有关Zinx框架的全局参数
 * 一些参数可以由用户通过zinx.json来配置
 */
type GlobalObj struct {
	TcpServer ziface.IServer //当前服务器zinx的全局Server对象
	//todo:存储全局Server对象有啥用?
	Host             string //服务器IP
	TcpPort          int    //服务器端口
	Name             string //服务器名称
	Version          string //服务器版本
	MaxPacketSize    uint32 //数据包的最大值
	MaxConn          int    //最大连接数
	WorkerPoolSize   uint32 //业务处理池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大存储数量
	ConfFilePath     string //配置文件的路径
	MaxMsgChanLen    int    //消息缓冲区大小
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {

	//读取配置文件zinx.json
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	//反序列化到GlobalObject中
	//todo:了解下json.Unmarshal的操作,为啥GlobalObject要用取地址
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "v0.4",
		TcpPort:          7777,
		Host:             "0.0.0.0",
		MaxConn:          120,
		MaxPacketSize:    4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		ConfFilePath:     "conf/zinx.json",
	}

	GlobalObject.Reload()
}
