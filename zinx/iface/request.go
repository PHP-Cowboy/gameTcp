package iface

type Request interface {
	GetConnection() Connection

	GetData() []byte
}
