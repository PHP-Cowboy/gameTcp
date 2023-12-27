package zNet

import (
	"fmt"
	"gameTcp/zinx/iface"
	"gameTcp/zinx/utils"
	"strconv"
)

type MsgHandler struct {
	Apis           map[uint32]iface.Router
	WorkerPoolSize uint32               //业务工作Worker池的数量
	TaskQueue      []chan iface.Request //Worker负责取任务的消息队列
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]iface.Router),
		WorkerPoolSize: utils.Global.WorkerPoolSize,
		TaskQueue:      make([]chan iface.Request, utils.Global.WorkerPoolSize),
	}
}

func (mh *MsgHandler) Do(req iface.Request) {
	handler, ok := mh.Apis[req.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", req.GetMsgId(), " is not FOUND!")
		return
	}

	//执行对应处理方法
	handler.PreHandle(req)
	handler.Handle(req)
	handler.AfterHandle(req)
}

func (mh *MsgHandler) AddRouter(msgId uint32, router iface.Router) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	//2 添加msg与api的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}

func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan iface.Request) {
	fmt.Println("Worker ID = ", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			mh.Do(request)
		}
	}
}

func (mh *MsgHandler) StartWorkerPool() {
	//遍历需要启动worker的数量，依此启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan iface.Request, utils.Global.MaxWorkerTaskLen)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandler) SendMsgToTaskQueue(req iface.Request) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := req.GetConnection().GetConnId() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", req.GetConnection().GetConnId(), " request msgID=", req.GetMsgId(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- req
}
