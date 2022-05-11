package demo

import (
	"github.com/NicholeGit/nade/framework"
	"github.com/pkg/errors"
)

func init() {
	err := framework.Register(&DemoProvider{})
	if err != nil {
		panic(errors.Wrap(err, "Register error"))
	}
}

type DemoProvider struct {
	framework.IServiceProvider

	c framework.IContainer
}

func (sp *DemoProvider) Name() string {
	return DemoKey
}

func (sp *DemoProvider) Register(_ framework.IContainer) framework.NewInstance {
	return NewService
}

func (sp *DemoProvider) IsDefer() bool {
	return false
}

func (sp *DemoProvider) Params(_ framework.IContainer) []interface{} {
	return []interface{}{sp.c}
}

func (sp *DemoProvider) Boot(c framework.IContainer) error {
	sp.c = c
	return nil
}
