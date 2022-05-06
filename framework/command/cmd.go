package command

import (
	"github.com/NicholeGit/nade/framework/util"
	"github.com/spf13/cobra"
)

// 初始化command相关命令
func initCmdCommand() *cobra.Command {
	cmdCommand.AddCommand(cmdListCommand)
	// cmdCommand.AddCommand(cmdCreateCommand)
	return cmdCommand
}

// 二级命令
var cmdCommand = &cobra.Command{
	Use:   "command",
	Short: "控制台命令相关",
	RunE: func(c *cobra.Command, args []string) error {
		_ = c.Help()
		return nil
	},
}

// cmdListCommand 列出所有的控制台命令
var cmdListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出所有控制台命令",
	RunE: func(c *cobra.Command, args []string) error {
		cmds := c.Root().Commands()
		ps := [][]string{}
		for _, cmd := range cmds {
			line := []string{cmd.Name(), cmd.Short}
			ps = append(ps, line)
		}
		util.PrettyPrint(ps)
		return nil
	},
}
