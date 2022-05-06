package console

import (
	"context"
	"fmt"

	"github.com/NicholeGit/nade/app/console/foo"
	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/command"
	"github.com/spf13/cobra"
)

// RunCommand  初始化根Command并运行
func RunCommand(container framework.IContainer) error {
	rootCmd, ctx := command.GetRootCmd(container)
	// 绑定业务的命令
	AddAppCommand(ctx, rootCmd)

	// 执行RootCommand
	return rootCmd.ExecuteContext(ctx)

}

// AddAppCommand 绑定业务的命令
func AddAppCommand(ctx context.Context, rootCmd *cobra.Command) {
	rootCmd.AddCommand(foo.FooCommand)

	// 每秒调用一次Foo命令
	c := ctx.Value(command.CommandCtxKey).(*command.CommandContextKey)
	err := c.AddCronCommand(ctx, "* * * * * *", foo.FooCommand)
	if err != nil {
		fmt.Println("AddCronCommand is error !", err)
	}
}
