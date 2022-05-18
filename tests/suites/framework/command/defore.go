package command

import (
	"context"
	"log"
	"sync/atomic"
	"testing"

	"github.com/NicholeGit/nade/framework/command"
	"github.com/NicholeGit/nade/framework/contract"
	_ "github.com/NicholeGit/nade/framework/provider"
	"github.com/NicholeGit/nade/tests"
)

func SetUp(t *testing.T) {
	log.Println("SetUp")
	gContainer = tests.InitBaseContainer()
	tServer := NewTestServer()
	var err error
	app := gContainer.MustMake(contract.AppKey).(contract.IApp)
	err = app.WithBaseFolder("../../../../")
	if err != nil {
		t.Error(err)
	}

	kernel := gContainer.MustMake(contract.KernelKey).(contract.IKernel)
	err = kernel.AddServers(tServer)
	if err != nil {
		t.Error(err)
	}
	gRootCmd, gCtx = command.GetRootCmd(gContainer)
}

func Before() {
	log.Println("Before")
}

type TestServer struct {
	contract.IServer
	isStart atomic.Value
}

func NewTestServer() *TestServer {
	tServer := &TestServer{}
	tServer.isStart.Store(false)
	return tServer
}

// Start 必须在 stop 中停止运行否则服务器关闭不掉
func (g *TestServer) Start(_ context.Context) error {
	log.Println("TestServer start")
	g.isStart.Store(true)
	for {
		if isRun := g.isStart.Load().(bool); !isRun {
			return nil
		}
		// server running !
	}
}

func (g *TestServer) Stop(_ context.Context) error {
	log.Println("TestServer stop")
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

func (g *TestServer) IsRun() bool {
	return g.isStart.Load().(bool)
}
