package contract

// AppKey 定义字符串凭证
const AppKey = "Nade:app"

// IApp 定义接口
type IApp interface {
	// AppID 表示当前这个app的唯一id, 可以用于分布式锁等
	AppID() string
	// Version 定义当前版本
	Version() string

	// BaseFolder 定义项目基础地址
	BaseFolder() string
	// ConfigFolder 定义了配置文件的路径
	ConfigFolder() string
	// LoadAppConfig 加载新的AppConfig，key为对应的函数转为小写下划线，比如ConfigFolder => config_folder
	LoadAppConfig(kv map[string]string)
	// RuntimeFolder 定义业务的运行中间态信息
	RuntimeFolder() string
	// LogFolder 定义了日志所在路径
	LogFolder() string
}
