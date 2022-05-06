package console

import (
	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/command"
	"github.com/spf13/cobra"
)

// RunCommand  初始化根Command并运行
func RunCommand(container framework.IContainer) error {
	rootCmd, ctx := command.GetRootCmd(container)
	// 绑定业务的命令
	AddAppCommand(rootCmd)

	// 执行RootCommand
	return rootCmd.ExecuteContext(ctx)

}

// 绑定业务的命令
func AddAppCommand(rootCmd *cobra.Command) {

}
