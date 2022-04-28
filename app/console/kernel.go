package console

import (
	"context"

	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/command"
	"github.com/spf13/cobra"
)

// RunCommand  初始化根Command并运行
func RunCommand(container framework.IContainer) error {

	rootCmd := command.GetRootCmd()
	// 绑定业务的命令
	AddAppCommand(rootCmd)

	ctx := context.Background()
	// 执行RootCommand
	err := rootCmd.ExecuteContext(context.WithValue(ctx, command.CommandCtxKey, command.CommandContextKey{
		Container: container,
	}))
	return err

}

// 绑定业务的命令
func AddAppCommand(rootCmd *cobra.Command) {

}
