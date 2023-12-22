package iface

type Pack interface {
	GetHeadLen() uint32             //获取包头长度方法
	Pack(Message) ([]byte, error)   //封包方法
	UnPack([]byte) (Message, error) //拆包方法
}
