package znet

import (
	"awesomeProject/utils"
	"awesomeProject/ziface"
	"fmt"
	"strconv"
)

type MsgHandle struct {
	Apis           map[uint32]ziface.IRounter //存放每个ID所对应的处理方法
	WorkerPoolSize uint32                     //业务工作Worker池大小
	TaskQuenue     []chan ziface.IRequest     //Worker负责取数据的消息队列
}

func (mh *MsgHandle) StartOneWorker(workId int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workId, " is Started...")
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandle) StartWorkerPool() {

	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQuenue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQuenue[i])
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(), " request msgID = ", request.GetMsgID(), " workerID = ", workerID)
	mh.TaskQuenue[workerID] <- request
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgid = ", request.GetMsgID(), " is not Found")
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)

	return
}

func (mh *MsgHandle) AddRounter(msgId uint32, rounter ziface.IRounter) {
	_, ok := mh.Apis[msgId]
	if ok {
		panic("repeated api, msgid = " + strconv.Itoa(int(msgId)))
	}

	mh.Apis[msgId] = rounter
	fmt.Println("Add api msgid = ", msgId)
	return
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRounter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQuenue:     make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}
