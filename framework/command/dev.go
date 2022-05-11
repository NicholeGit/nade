package command

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/contract"
	"github.com/NicholeGit/nade/framework/util"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// devConfig 代表调试模式的配置信息
type devConfig struct {
	Port    string   // 调试模式最终监听的端口，默认为8070
	Backend struct { // 后端调试模式配置
		RefreshTime   int    // 调试模式后端更新时间，如果文件变更，等待3s才进行一次更新，能让频繁保存变更更为顺畅, 默认1s
		Port          string // 后端监听端口， 默认 8072
		MonitorFolder string // 监听文件夹，默认为AppFolder
	}
}

// 初始化配置文件
func initDevConfig(c framework.IContainer) *devConfig {
	// 设置默认值
	devConfig := &devConfig{
		Port: "8087",
		Backend: struct {
			RefreshTime   int
			Port          string
			MonitorFolder string
		}{
			1,
			"8072",
			"",
		},
	}
	// 容器中获取配置服务
	configer := c.MustMake(contract.ConfigKey).(contract.IConfig)
	// 每个配置项进行检查
	if configer.IsExist("app.dev.port") {
		devConfig.Port = configer.GetString("app.dev.port")
	}
	if configer.IsExist("app.dev.backend.refresh_time") {
		devConfig.Backend.RefreshTime = configer.GetInt("app.dev.backend.refresh_time")
	}
	if configer.IsExist("app.dev.backend.port") {
		devConfig.Backend.Port = configer.GetString("app.dev.backend.port")
	}
	// monitorFolder 默认使用目录服务的AppFolder()
	monitorFolder := configer.GetString("app.dev.backend.monitor_folder")
	if monitorFolder == "" {
		appService := c.MustMake(contract.AppKey).(contract.IApp)
		devConfig.Backend.MonitorFolder = appService.AppFolder()
	} else {
		devConfig.Backend.MonitorFolder = monitorFolder
	}

	return devConfig
}

// Proxy 代表serve启动的服务器代理
type Proxy struct {
	devConfig  *devConfig // 配置文件
	backendPid int        // 当前的backend服务的pid
}

// NewProxy 初始化一个Proxy
func NewProxy(c framework.IContainer) *Proxy {
	devConfig := initDevConfig(c)
	return &Proxy{
		devConfig: devConfig,
	}
}

// 重新启动一个 proxy 网关
func (p *Proxy) newProxyReverseProxy(frontend, backend *url.URL) *httputil.ReverseProxy {
	if p.backendPid == 0 {
		fmt.Println("后端服务都不存在")
		return nil
	}

	// 后端服务存在
	if p.backendPid != 0 {
		return httputil.NewSingleHostReverseProxy(backend)
	}

	// 两个都有进程
	// 先创建一个后端服务的directory
	director := func(req *http.Request) {
		req.URL.Scheme = backend.Scheme
		req.URL.Host = backend.Host
	}

	// 定义一个NotFoundErr
	NotFoundErr := errors.New("response is 404, need to redirect")
	return &httputil.ReverseProxy{
		Director: director, // 先转发到后端服务
		ModifyResponse: func(response *http.Response) error {
			// 如果后端服务返回了404，我们返回NotFoundErr 会进入到errorHandler中
			if response.StatusCode == 404 {
				return NotFoundErr
			}
			return nil
		},
		ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
			// 判断 Error 是否为NotFoundError, 是的话则进行前端服务的转发，重新修改writer
			if errors.Is(err, NotFoundErr) {
				httputil.NewSingleHostReverseProxy(frontend).ServeHTTP(writer, request)
			}
		}}
}

// rebuildBackend 重新编译后端
func (p *Proxy) rebuildBackend() error {
	// 重新编译 Nade
	bin := os.Args[0]
	cmdBuild := exec.Command(bin, "build", "backend")
	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr
	if err := cmdBuild.Start(); err == nil {
		err = cmdBuild.Wait()
		if err != nil {
			return err
		}
	}
	return nil
}

// restartBackend 启动后端服务
func (p *Proxy) restartBackend() error {

	// 杀死之前的进程
	if p.backendPid != 0 {
		err := syscall.Kill(p.backendPid, syscall.SIGKILL)
		if err != nil {
			return errors.Wrap(err, "restartBackend syscall.kill is error")
		}
		p.backendPid = 0
	}

	// 设置随机端口，真实后端的端口
	port := p.devConfig.Backend.Port
	nadeAddress := fmt.Sprintf(":" + port)

	bin := os.Args[0]
	// 使用命令行启动后端进程
	cmd := exec.Command(bin, "app", "start", "--address="+nadeAddress)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("启动后端服务: ", "http://127.0.0.1:"+port)
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	p.backendPid = cmd.Process.Pid
	fmt.Println("后端服务pid:", p.backendPid)
	return nil
}

// 启动proxy服务，并且根据参数启动前端服务或者后端服务
func (p *Proxy) startProxy(startBackend bool) error {
	var backendURL, frontendURL *url.URL
	var err error

	// 启动后端
	if startBackend {
		if err := p.restartBackend(); err != nil {
			return err
		}
	}

	if backendURL, err = url.Parse(fmt.Sprintf("%s%s", "http://127.0.0.1:", p.devConfig.Backend.Port)); err != nil {
		return err
	}

	// 设置反向代理
	proxyReverse := p.newProxyReverseProxy(frontendURL, backendURL)
	proxyServer := &http.Server{
		Addr:    "127.0.0.1:" + p.devConfig.Port,
		Handler: proxyReverse,
	}

	fmt.Println("代理服务启动:", "http://"+proxyServer.Addr)
	// 启动proxy服务
	err = proxyServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

// monitorBackend 监听应用文件
func (p *Proxy) monitorBackend() error {
	// 监听
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// 开启监听目标文件夹
	appFolder := p.devConfig.Backend.MonitorFolder
	fmt.Println("监控文件夹：", appFolder)
	// 监听所有子目录，需要使用filepath.walk
	err = filepath.Walk(appFolder, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			return nil
		}
		// 如果是隐藏的目录比如 . 或者 .. 则不用进行监控
		if util.IsHiddenDirectory(path) {
			return nil
		}
		return watcher.Add(path)
	})
	if err != nil {
		return errors.Wrap(err, "monitorBackend filepath.Walk is error")
	}

	// 开启计时时间机制
	refreshTime := p.devConfig.Backend.RefreshTime
	t := time.NewTimer(time.Duration(refreshTime) * time.Second)
	// 先停止计时器
	t.Stop()
	for {
		select {
		case <-t.C:
			// 计时器时间到了，代表之前有文件更新事件重置过计时器
			// 即有文件更新
			fmt.Println("...检测到文件更新，重启服务开始...")
			fmt.Println("...期间请不要发送任何请求...")
			if err := p.rebuildBackend(); err != nil {
				fmt.Println("重新编译失败：", err.Error())
			} else {
				if err := p.restartBackend(); err != nil {
					fmt.Println("重新启动失败：", err.Error())
				}
			}
			fmt.Println("...检测到文件更新，重启服务结束...")
			// 停止计时器
			t.Stop()
		case ev, ok := <-watcher.Events:
			fmt.Println("文件变化", ev.Name)
			if !ok {
				continue
			}
			// 有文件更新事件，重置计时器
			t.Reset(time.Duration(refreshTime) * time.Second)
		case err, ok := <-watcher.Errors:
			if !ok {
				continue
			}
			// 如果有文件监听错误，则停止计时器
			fmt.Println("监听文件夹错误：", err.Error())
			t.Reset(time.Duration(refreshTime) * time.Second)
		}
	}
}

// 初始化Dev命令
func initDevCommand() *cobra.Command {
	devCommand.AddCommand(devBackendCommand)
	return devCommand
}

// devCommand 为调试模式的一级命令
var devCommand = &cobra.Command{
	Use:   "dev",
	Short: "调试模式",
	RunE: func(c *cobra.Command, args []string) error {
		return c.Help()
	},
}

// devBackendCommand 启动后端调试模式
var devBackendCommand = &cobra.Command{
	Use:   "backend",
	Short: "启动后端调试模式",
	RunE: func(c *cobra.Command, args []string) error {
		container := GetCommandContextKey(c).Container()
		proxy := NewProxy(container)
		go proxy.monitorBackend()
		if err := proxy.startProxy(true); err != nil {
			return err
		}
		return nil
	},
}
