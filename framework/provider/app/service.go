package app

import (
	"flag"
	"path/filepath"

	"github.com/NicholeGit/nade/framework/command"
	"github.com/NicholeGit/nade/framework/util"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/NicholeGit/nade/framework"
)

// NadeApp 代表nade框架的App实现
type NadeApp struct {
	container  framework.IContainer // 服务容器
	baseFolder string               // 基础路径
	appId      string               // 表示当前这个app的唯一id, 可以用于分布式锁等

	configMap map[string]string // 配置加载
}

// NewNadeApp 初始化nadeApp
func NewNadeApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}
	// 有两个参数，一个是容器，一个是baseFolder
	container := params[0].(framework.IContainer)
	baseFolder := params[1].(string)
	if baseFolder == "" {
		// 如果没有设置，则使用参数
		baseFolder = command.BaseFolder
		flag.StringVar(&baseFolder, "base_folder", "", "base_folder参数, 默认为当前路径")
		flag.Parse()
	}

	appId := uuid.New().String()
	configMap := map[string]string{}
	return &NadeApp{baseFolder: baseFolder, container: container, appId: appId, configMap: configMap}, nil
}

// AppID 表示这个App的唯一ID
func (app NadeApp) AppID() string {
	return app.appId
}

// Version 实现版本
func (app NadeApp) Version() string {
	return "0.0.3"
}

// BaseFolder 表示基础目录，可以代表开发场景的目录，也可以代表运行时候的目录
func (app NadeApp) BaseFolder() string {
	if app.baseFolder != "" {
		return app.baseFolder
	}

	// 如果参数也没有，使用默认的当前路径
	return util.GetExecDirectory()
}

// ConfigFolder  表示配置文件地址
func (app NadeApp) ConfigFolder() string {
	if val, ok := app.configMap["config_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "config")
}

func (app *NadeApp) StorageFolder() string {
	if val, ok := app.configMap["storage_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "storage")
}

// RuntimeFolder 定义业务的运行中间态信息
func (app NadeApp) RuntimeFolder() string {
	if val, ok := app.configMap["runtime_folder"]; ok {
		return val
	}
	return filepath.Join(app.StorageFolder(), "runtime")
}

// LogFolder 表示日志存放地址
func (app NadeApp) LogFolder() string {
	if val, ok := app.configMap["log_folder"]; ok {
		return val
	}
	return filepath.Join(app.StorageFolder(), "log")
}

// LoadAppConfig 加载配置map
func (app *NadeApp) LoadAppConfig(kv map[string]string) {
	for key, val := range kv {
		app.configMap[key] = val
	}
}

// AppFolder 代表app目录
func (app *NadeApp) AppFolder() string {
	if val, ok := app.configMap["app_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app")
}

// MiddlewareFolder 定义业务自己定义的中间件
func (app NadeApp) MiddlewareFolder() string {
	if val, ok := app.configMap["middleware_folder"]; ok {
		return val
	}
	return filepath.Join(app.HttpFolder(), "middleware")
}

func (app NadeApp) HttpFolder() string {
	if val, ok := app.configMap["http_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app", "http")
}
