package http

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/NicholeGit/nade/app/http/module/demo"
	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/contract"
)

// Routes 绑定业务层路由
func Routes(container framework.IContainer, r *gin.Engine) {
	configService := container.MustMake(contract.ConfigKey).(contract.IConfig)
	//
	// // /路径先去./dist目录下查找文件是否存在，找到使用文件服务提供服务
	// r.Use(static.Serve("/", static.LocalFile("./dist", false)))
	//
	// 如果配置了swagger，则显示swagger的中间件
	if configService.GetBool("app.swagger") == true {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 动态路由定义
	demo.Register(r)
}
