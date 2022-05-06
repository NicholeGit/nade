package command

import (
	"github.com/NicholeGit/nade/framework"
	"github.com/spf13/cobra"
)

// 初始化provider相关服务
func initProviderCommand() *cobra.Command {
	providerCommand.AddCommand(providerListCommand)
	return providerCommand
}

// providerCommand 二级命令
var providerCommand = &cobra.Command{
	Use:   "provider",
	Short: "服务提供相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			_ = c.Help()
		}
		return nil
	},
}

// providerListCommand 列出容器内的所有服务
var providerListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出容器内的所有服务",
	RunE: func(c *cobra.Command, args []string) error {
		container := GetCommandContextKey(c).Container()
		hadeContainer := container.(*framework.NadeContainer)
		// 获取字符串凭证
		list := hadeContainer.NameList()
		// 打印
		for _, line := range list {
			println(line)
		}
		return nil
	},
}
