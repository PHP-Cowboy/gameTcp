package zNet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"gameTcp/zinx/iface"
	"gameTcp/zinx/utils"
)

type Pack struct{}

// 封包拆包实例初始化方法
func NewPack() *Pack {
	return &Pack{}
}

// 获取包头长度方法
func (p *Pack) GetHeadLen() uint32 {
	return 8
}

// 封包方法
func (p *Pack) Pack(msg iface.Message) (data []byte, err error) {
	buffer := bytes.NewBuffer([]byte{})

	//写dataLen
	if err = binary.Write(buffer, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return
	}

	//写msgID
	if err = binary.Write(buffer, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return
	}

	//写data数据
	if err = binary.Write(buffer, binary.LittleEndian, msg.GetData()); err != nil {
		return
	}

	data = buffer.Bytes()

	return
}

// 拆包方法
func (p *Pack) UnPack(data []byte) (iface.Message, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(data)

	msg := &Msg{}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Data); err != nil {
		return nil, err
	}

	if utils.Global.MaxPacketSize > 0 && utils.Global.MaxPacketSize < msg.DataLen {
		return nil, errors.New("Too large msg data recieved")
	}

	return msg, nil
}
