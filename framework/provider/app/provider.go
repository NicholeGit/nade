package app

import (
	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/contract"
)

// NadeAppProvider 提供App的具体实现方法
type NadeAppProvider struct {
	BaseFolder string
}

// Register 注册nadeApp方法
func (n *NadeAppProvider) Register(container framework.IContainer) framework.NewInstance {
	return NewNadeApp
}

// Boot 启动调用
func (n *NadeAppProvider) Boot(container framework.IContainer) error {
	return nil
}

// IsDefer 是否延迟初始化
func (n *NadeAppProvider) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (n *NadeAppProvider) Params(container framework.IContainer) []interface{} {
	return []interface{}{container, n.BaseFolder}
}

// Name 获取字符串凭证
func (n *NadeAppProvider) Name() string {
	return contract.AppKey
}
