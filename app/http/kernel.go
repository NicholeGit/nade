package http

import (
	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/contract"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// NewGinEngine 创建了一个绑定了路由的Web引擎
func NewGinEngine(container framework.IContainer) (*gin.Engine, error) {
	// 设置为Release，为的是默认在启动中不输出调试信息
	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	// 默认启动一个Web引擎
	r := gin.New()

	// 默认注册recovery中间件
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// 业务绑定路由操作
	Routes(container, r)
	// 返回绑定路由后的Web引擎
	return r, nil
}

func InitServer(container framework.IContainer) error {
	// 创建 gin server
	ginEngine, err := NewGinEngine(container)
	if err != nil {
		return errors.Wrap(err, "NewHttpEngine is error !")
	}

	var appAddress string
	envService := container.MustMake(contract.EnvKey).(contract.IEnv)
	if envService.Get("ADDRESS") != "" {
		appAddress = envService.Get("ADDRESS")
	} else {
		configService := container.MustMake(contract.ConfigKey).(contract.IConfig)
		if configService.IsExist("app.address") {
			appAddress = configService.GetString("app.address")
		} else {
			appAddress = ":8888"
		}
	}
	ginServer := NewGinServer(ginEngine, appAddress)

	// 注册所有 server 到 Kernel 组件中
	kernel := container.MustMake(contract.KernelKey).(contract.IKernel)
	err = kernel.AddServers(ginServer)
	if err != nil {
		return errors.Wrap(err, "InitServer  AddServers ginServer")
	}
	return nil
}
