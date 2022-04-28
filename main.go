package main

import (
	"fmt"

	"github.com/NicholeGit/nade/app/console"
	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/provider/app"
	"github.com/NicholeGit/nade/framework/provider/config"
	"github.com/NicholeGit/nade/framework/provider/env"
)

func main() {
	// 初始化服务容器
	container := framework.NewNadeContainer()

	// 绑定App服务提供者
	container.Bind(&app.NadeAppProvider{})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&env.NadeEnvProvider{})
	container.Bind(&config.NadeConfigProvider{})

	// 运行root命令
	err := console.RunCommand(container)
	if err != nil {
		fmt.Println("RunCommand is err !", err)
	}

}
