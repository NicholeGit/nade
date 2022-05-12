package tests

import (
	"github.com/NicholeGit/nade/framework"
	// init provider
	_ "github.com/NicholeGit/nade/framework/provider"
	"github.com/pkg/errors"
)

func InitBaseContainer() framework.IContainer {
	// 初始化服务容器
	container := framework.NewNadeContainer()

	if err := container.BindAll(); err != nil {
		panic(errors.Wrap(err, "bindAll is err"))
	}
	return container
}
