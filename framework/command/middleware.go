package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/NicholeGit/nade/framework/util"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/NicholeGit/nade/framework/contract"
)

// 初始化中间件相关命令
func initMiddlewareCommand() *cobra.Command {
	middlewareCommand.AddCommand(middlewareAllCommand)
	middlewareCommand.AddCommand(middlewareCreateCommand)
	return middlewareCommand
}

// middlewareCommand 中间件二级命令
var middlewareCommand = &cobra.Command{
	Use:   "middleware",
	Short: "中间件相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

// middlewareAllCommand 显示所有安装的中间件
var middlewareAllCommand = &cobra.Command{
	Use:   "list",
	Short: "显示所有中间件",
	RunE: func(c *cobra.Command, args []string) error {
		container := GetCommandContextKey(c).Container()
		appService := container.MustMake(contract.AppKey).(contract.IApp)

		middlewarePath := path.Join(appService.BaseFolder(), "app", "http", "middleware")

		// 读取文件夹
		files, err := ioutil.ReadDir(middlewarePath)
		if err != nil {
			return err
		}
		for _, f := range files {
			if f.IsDir() {
				fmt.Println(f.Name())
			}
		}
		return nil
	},
}

// providerCreateCommand 创建一个新的服务，包括服务提供者，服务接口协议，服务实例
var middlewareCreateCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Short:   "创建一个中间件",
	RunE: func(c *cobra.Command, args []string) error {
		container := GetCommandContextKey(c).Container()
		fmt.Println("创建一个中间件")
		var name string
		var folder string
		{
			prompt := &survey.Input{
				Message: "请输入中间件名称：",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}
		{
			prompt := &survey.Input{
				Message: "请输入中间件所在目录名称(默认: 同中间件名称):",
			}
			err := survey.AskOne(prompt, &folder)
			if err != nil {
				return err
			}
		}
		if folder == "" {
			folder = name
		}

		// 判断是否有重复中间件
		app := container.MustMake(contract.AppKey).(contract.IApp)
		pFolder := app.MiddlewareFolder()
		if util.Exists(pFolder) {
			subFolders, err := util.SubDir(pFolder)
			if err != nil {
				return err
			}
			for _, s := range subFolders {
				if s == folder {
					fmt.Println("目录已经存在")
					return nil
				}
			}
		}

		// 开始创建文件
		if err := os.MkdirAll(filepath.Join(pFolder, folder), 0700); err != nil {
			return err
		}
		funcs := template.FuncMap{"title": strings.Title}
		{
			//  创建
			file := filepath.Join(pFolder, folder, "middleware.go")
			f, err := os.Create(file)
			if err != nil {
				return errors.Cause(err)
			}

			t := template.Must(template.New("middleware").Funcs(funcs).Parse(middlewareTmp))
			if err := t.Execute(f, name); err != nil {
				return errors.Cause(err)
			}
		}
		fmt.Println("创建中间件成功, 文件夹地址:", filepath.Join(pFolder, folder))
		return nil
	},
}

var middlewareTmp = `package {{.}}

import "github.com/gohade/hade/framework/gin"

// {{.|title}}Middleware 代表中间件函数
func {{.|title}}Middleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()
	}
}

`
