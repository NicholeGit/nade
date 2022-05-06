package command

import (
	"context"
	"sync"

	"github.com/NicholeGit/nade/framework"
	"github.com/spf13/cobra"
)

var once sync.Once

const CommandCtxKey = "CommandCtx"

type CommandContextKey struct {
	Container framework.IContainer
}

var (
	BaseFolder string
)

func init() {
	gRootCmd.PersistentFlags().StringVar(&BaseFolder, "base_folder", ".", "base_folder参数, 默认为当前路径")
}

var (
	gRootCmd = &cobra.Command{
		// 定义根命令的关键字
		Use: "nade",
		// 简短介绍
		Short: "nade 命令",
		// 根命令的详细介绍
		Long: "nade 框架提供的命令行工具，使用这个命令行工具能很方便执行框架自带命令，也能很方便编写业务命令",
		// 根命令的执行函数
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
		// 不需要出现 cobra 默认的 completion 子命令
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	gCtx = context.Background()
)

func GetRootCmd(container framework.IContainer) (*cobra.Command, context.Context) {
	once.Do(func() {
		// 绑定框架的命令
		AddKernelCommands(gRootCmd)
		gCtx = context.WithValue(gCtx, CommandCtxKey, CommandContextKey{
			Container: container,
		})
	})
	return gRootCmd, gCtx
}

// GetCommandContext 从 cobra.Command 获取 CommandContext
func GetCommandContext(cmd *cobra.Command) CommandContextKey {
	return cmd.Context().Value(CommandCtxKey).(CommandContextKey)
}
