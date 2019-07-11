package ziface

/*
 * 消息管理抽象层
 */
type IMsgHandle interface {

	//以非阻塞的方式处理消息
	DoMsgHandler(request IRequest)

	//为消息添加具体的处理逻辑
	AddRounter(msgId uint32, rounter IRounter)

	StartWorkerPool()

	SendMsgToTaskQueue(request IRequest)
}
