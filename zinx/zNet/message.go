package zNet

type Msg struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

func NewMsg(msgId uint32, data []byte) *Msg {
	return &Msg{
		Id:      msgId,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

func (m *Msg) GetMsgId() uint32 {
	return m.Id
}

func (m *Msg) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Msg) GetData() []byte {
	return m.Data
}

func (m *Msg) SetMsgId(msgId uint32) {
	m.Id = msgId
}
func (m *Msg) SetDataLen(len uint32) {
	m.DataLen = len
}
func (m *Msg) SetData(data []byte) {
	m.Data = data
}
