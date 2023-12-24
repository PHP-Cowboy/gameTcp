package zNet

import (
	"fmt"
	"gameTcp/zinx/iface"
	"strconv"
)

type MsgHandler struct {
	Apis map[uint32]iface.Router
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{Apis: make(map[uint32]iface.Router)}
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
