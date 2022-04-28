package command

import (
	"fmt"

	"github.com/NicholeGit/nade/framework/contract"
	"github.com/NicholeGit/nade/framework/util"
	"github.com/spf13/cobra"
)

// initEnvCommand 获取env相关的命令
func initEnvCommand() *cobra.Command {
	envCommand.AddCommand(envListCommand)
	return envCommand
}

// envCommand 获取当前的App环境
var envCommand = &cobra.Command{
	Use:   "env",
	Short: "获取当前的App环境",
	Run: func(c *cobra.Command, args []string) {

		// 获取env环境
		container := GetCommandContext(c).Container
		envService := container.MustMake(contract.EnvKey).(contract.IEnv)
		// 打印环境
		fmt.Println("environment:", envService.AppEnv())
	},
}

// envListCommand 获取所有的App环境变量
var envListCommand = &cobra.Command{
	Use:   "list",
	Short: "获取所有的环境变量",
	Run: func(c *cobra.Command, args []string) {
		// 获取env环境
		container := GetCommandContext(c).Container
		envService := container.MustMake(contract.EnvKey).(contract.IEnv)
		envs := envService.All()
		var outs [][]string
		for k, v := range envs {
			outs = append(outs, []string{k, v})
		}
		util.PrettyPrint(outs)
	},
}
