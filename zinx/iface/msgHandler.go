package iface

type MsgHandler interface {
	Do(request Request)
	AddRouter(msgId uint32, router Router)
	StartWorkerPool()
	SendMsgToTaskQueue(request Request)
}
