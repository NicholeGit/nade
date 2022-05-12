package id

import (
	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/contract"
	"github.com/pkg/errors"
)

func init() {
	err := framework.Register(&NadeIDProvider{})
	if err != nil {
		panic(errors.Wrap(err, "Register error"))
	}
}

type NadeIDProvider struct{}

// Register registe a new function for make a service instance
func (provider *NadeIDProvider) Register(_ framework.IContainer) framework.NewInstance {
	return NewNadeIDService
}

// Boot will called when the service instantiate
func (provider *NadeIDProvider) Boot(_ framework.IContainer) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *NadeIDProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *NadeIDProvider) Params(_ framework.IContainer) []interface{} {
	return []interface{}{}
}

// Name define the name for this service
func (provider *NadeIDProvider) Name() string {
	return contract.IDKey
}
