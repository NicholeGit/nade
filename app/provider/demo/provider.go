package demo

import (
	"github.com/NicholeGit/nade/framework"
	"github.com/pkg/errors"
)

func init() {
	err := framework.Register(&Provider{})
	if err != nil {
		panic(errors.Wrap(err, "Register error"))
	}
}

type Provider struct {
	framework.IServiceProvider

	c framework.IContainer
}

func (sp *Provider) Name() string {
	return DemoKey
}

func (sp *Provider) Register(_ framework.IContainer) framework.NewInstance {
	return NewService
}

func (sp *Provider) IsDefer() bool {
	return false
}

func (sp *Provider) Params(_ framework.IContainer) []interface{} {
	return []interface{}{sp.c}
}

func (sp *Provider) Boot(c framework.IContainer) error {
	sp.c = c
	return nil
}
