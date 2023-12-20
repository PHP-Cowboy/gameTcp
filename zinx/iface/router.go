package iface

type Router interface {
	PreHandle(request Request)
	Handle(request Request)
	AfterHandle(request Request)
}
