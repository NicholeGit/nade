package main

import (
	"fmt"

	"github.com/NicholeGit/nade/app/console"
	"github.com/NicholeGit/nade/framework"
	_ "github.com/NicholeGit/nade/framework/provider"
	"github.com/pkg/errors"
)

func main() {
	// 初始化服务容器
	container := framework.NewNadeContainer()

	err := container.BindAll()
	if err != nil {
		panic(errors.Wrap(err, "BindAll is err !"))
	}

	// 运行root命令
	err = console.RunCommand(container)
	if err != nil {
		fmt.Println("RunCommand is err !", err)
	}

}
