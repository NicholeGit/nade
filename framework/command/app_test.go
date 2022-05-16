package command_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/NicholeGit/nade/framework/command"
	"github.com/NicholeGit/nade/framework/contract"
	"github.com/NicholeGit/nade/tests"
)

func TestAppStart(t *testing.T) {
	container := tests.InitBaseContainer()
	tServer := &TestServer{}
	var err error
	rootCmd, ctx := command.GetRootCmd(container)
	app := container.MustMake(contract.AppKey).(contract.IApp)
	app.WithBaseFolder("../../")

	kernel := container.MustMake(contract.KernelKey).(contract.IKernel)
	err = kernel.AddServers(tServer)
	if err != nil {
		t.Error(err)
	}

	go func() {
		rootCmd.SetArgs([]string{"app", "start"})
		rootCmd, err = rootCmd.ExecuteContextC(ctx)
		if err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(time.Second)
	rootCmd.ResetFlags()
	rootCmd.SetArgs([]string{"app", "stop"})
	_, err = rootCmd.ExecuteContextC(ctx)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second)
}

type TestServer struct {
	contract.IServer
	isStart atomic.Value
}

// Start 必须在 stop 中停止运行否则服务器关闭不掉
func (g *TestServer) Start(_ context.Context) error {
	fmt.Println("TestServer start")
	g.isStart.Store(true)
	for {
		if isRun := g.isStart.Load().(bool); !isRun {
			return nil
		}
		// server running !
	}
}

func (g *TestServer) Stop(_ context.Context) error {
	fmt.Println("TestServer stop")
	// 跳出 start 中的 loop
	g.isStart.Store(false)
	return nil
}

func (g *TestServer) Name() string {
	return "TestServer"
}

func (g *TestServer) Addr() string {
	return ""
}
