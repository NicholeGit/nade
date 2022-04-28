package config

import (
	"path/filepath"

	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/contract"
)

type NadeConfigProvider struct{}

// Register a new function for make a service instance
func (provider *NadeConfigProvider) Register(c framework.IContainer) framework.NewInstance {
	return NewNadeConfig
}

// Boot will call when the service instantiate
func (provider *NadeConfigProvider) Boot(c framework.IContainer) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *NadeConfigProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *NadeConfigProvider) Params(c framework.IContainer) []interface{} {
	appService := c.MustMake(contract.AppKey).(contract.App)
	envService := c.MustMake(contract.EnvKey).(contract.IEnv)
	env := envService.AppEnv()
	// 配置文件夹地址
	configFolder := appService.ConfigFolder()
	envFolder := filepath.Join(configFolder, env)
	return []interface{}{c, envFolder, envService.All()}
}

// Name define the name for this service
func (provider *NadeConfigProvider) Name() string {
	return contract.ConfigKey
}
