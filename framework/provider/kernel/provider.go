package kernel

import (
	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/contract"
	"github.com/pkg/errors"
)

func init() {
	err := framework.Register(&NadeKernelProvider{})
	if err != nil {
		panic(errors.Wrap(err, "Register error"))
	}
}

// NadeKernelProvider 提供web引擎
type NadeKernelProvider struct{}

// Register 注册服务提供者
func (provider *NadeKernelProvider) Register(_ framework.IContainer) framework.NewInstance {
	return NewNadeKernelService
}

// Boot 启动的时候判断是否由外界注入了Engine，如果注入的化，用注入的，如果没有，重新实例化
func (provider *NadeKernelProvider) Boot(_ framework.IContainer) error {
	return nil
}

// IsDefer 引擎的初始化我们希望开始就进行初始化
func (provider *NadeKernelProvider) IsDefer() bool {
	return false
}

// Params 参数就是一个HttpEngine
func (provider *NadeKernelProvider) Params(c framework.IContainer) []interface{} {
	return []interface{}{c}
}

// Name 提供凭证
func (provider *NadeKernelProvider) Name() string {
	return contract.KernelKey
}

func (provider NadeKernelProvider) DependOn() []string {
	return []string{contract.ConfigKey, contract.EnvKey}
}
