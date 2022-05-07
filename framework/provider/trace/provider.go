package trace

import (
	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/contract"
	"github.com/pkg/errors"
)

func init() {
	err := framework.Register(&NadeTraceProvider{})
	if err != nil {
		panic(errors.Wrap(err, "Register error"))
	}
}

type NadeTraceProvider struct {
	c framework.IContainer
}

// Register registe a new function for make a service instance
func (provider *NadeTraceProvider) Register(_ framework.IContainer) framework.NewInstance {
	return NewNadeTraceService
}

// Boot will called when the service instantiate
func (provider *NadeTraceProvider) Boot(c framework.IContainer) error {
	provider.c = c
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *NadeTraceProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *NadeTraceProvider) Params(_ framework.IContainer) []interface{} {
	return []interface{}{provider.c}
}

// Name define the name for this service
func (provider *NadeTraceProvider) Name() string {
	return contract.TraceKey
}

func (provider *NadeTraceProvider) DependOn() []string {
	return []string{contract.IDKey}
}
