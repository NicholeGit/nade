package foo

import (
	"github.com/NicholeGit/nade/framework/command"
	"github.com/NicholeGit/nade/framework/contract"
	"github.com/spf13/cobra"
)

var FooCommand = &cobra.Command{
	Use:   "foo",
	Short: "foo",
	RunE: func(c *cobra.Command, args []string) error {
		container := command.GetCommandContextKey(c).Container()
		logger := container.MustMake(contract.LogKey).(contract.Log)

		logger.Debug(nil, "this is foo command", nil)
		return nil
	},
}
