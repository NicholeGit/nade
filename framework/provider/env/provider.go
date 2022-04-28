package env

import (
	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/contract"
)

type NadeEnvProvider struct {
	Folder string
}

// Register registe a new function for make a service instance
func (provider *NadeEnvProvider) Register(c framework.IContainer) framework.NewInstance {
	return NewNadeEnv
}

// Boot will called when the service instantiate
func (provider *NadeEnvProvider) Boot(c framework.IContainer) error {
	app := c.MustMake(contract.AppKey).(contract.App)
	provider.Folder = app.BaseFolder()
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *NadeEnvProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *NadeEnvProvider) Params(c framework.IContainer) []interface{} {
	return []interface{}{provider.Folder}
}

// Name define the name for this service
func (provider *NadeEnvProvider) Name() string {
	return contract.EnvKey
}
