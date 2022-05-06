package command

import (
	"context"
	"log"
	"sync"

	"github.com/NicholeGit/nade/framework"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

const CommandCtxKey = "CommandCtx"

// CronSpec 保存Cron命令的信息，用于展示
type CronSpec struct {
	Type        string
	Cmd         *cobra.Command
	Spec        string
	ServiceName string
	Id          int
}

type CommandContextKey struct {
	container framework.IContainer

	// Command支持cron，只在RootCommand中有这个值
	cron *cron.Cron
	// 对应Cron命令的说明文档
	cronSpecs []CronSpec
}

func init() {
	gRootCmd.PersistentFlags().StringVar(&BaseFolder, "base_folder", ".", "base_folder参数, 默认为当前路径")
}

var (
	BaseFolder string
	once       sync.Once
	gRootCmd   = &cobra.Command{
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
		gCtx = context.WithValue(gCtx, CommandCtxKey, &CommandContextKey{
			container: container,
			cron:      cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor))),
			cronSpecs: []CronSpec{},
		})
	})
	return gRootCmd, gCtx
}

// GetCommandContextKey 从 cobra.Command 获取 CommandContext
func GetCommandContextKey(cmd *cobra.Command) *CommandContextKey {
	return cmd.Context().Value(CommandCtxKey).(*CommandContextKey)
}

func (k *CommandContextKey) Container() framework.IContainer {
	return k.container
}

func (k *CommandContextKey) CronSpecs() []CronSpec {
	return k.cronSpecs
}
func (k *CommandContextKey) Cron() *cron.Cron {
	return k.cron
}

func (k *CommandContextKey) AddCronCommand(ctx context.Context, spec string, cmd *cobra.Command) error {
	// 增加调用函数
	id, err := k.cron.AddFunc(spec, func() {
		// 制作一个rootCommand，必须放在这个里面做复制，否则会产生竞态
		var cronCmd cobra.Command
		cronCmd = *cmd
		cronCmd.ResetCommands()
		cronCmd.SetArgs([]string{})
		// cronCmd.SetParantNull()
		// cronCmd.SetContainer(root.GetContainer())

		// 如果后续的command出现panic，这里要捕获
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		err := cronCmd.ExecuteContext(ctx)
		if err != nil {
			// 打印出err信息
			log.Println(err)
		}
	})
	if err != nil {
		return err
	}
	// 增加说明信息
	k.cronSpecs = append(k.cronSpecs, CronSpec{
		Type: "normal-cron",
		Cmd:  cmd,
		Spec: spec,
		Id:   int(id),
	})
	return nil

}
