package tests

import (
	"github.com/NicholeGit/nade/framework"
	_ "github.com/NicholeGit/nade/framework/provider"
	"github.com/pkg/errors"
)

func InitBaseContainer() framework.IContainer {
	// 初始化服务容器
	container := framework.NewNadeContainer()
	err := container.BindAll()
	if err != nil {
		panic(errors.Wrap(err, "BindAll is err !"))
	}
	return container
}
