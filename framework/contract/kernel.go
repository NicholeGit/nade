package contract

import (
	"context"
)

const KernelKey = "nade:kernel"

// IServer is server.
type IServer interface {
	Start(context.Context) error
	Stop(context.Context) error
	Name() string
	WithAddr(string) error
	Addr() string
}

// IKernel 接口提供框架最核心的结构
type IKernel interface {
	// Servers 得到所有server
	Servers() []IServer
	// AddServers 注册server
	AddServers(srv ...IServer) error
	// GetServer 根据name得到server
	GetServer(string) IServer
	// Info 服务器信息
	Info() string
}
