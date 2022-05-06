package command

import "github.com/spf13/cobra"

// AddKernelCommands will add all command/* to root command
func AddKernelCommands(root *cobra.Command) {
	// app 命令
	root.AddCommand(initAppCommand())
	// build 命令
	root.AddCommand(initBuildCommand())
	// go build
	root.AddCommand(goCommand)
	// cmd
	root.AddCommand(initCmdCommand())
	// env 命令
	root.AddCommand(initEnvCommand())
	// config 命令
	root.AddCommand(initConfigCommand())
}
