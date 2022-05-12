package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NicholeGit/nade/framework/contract"
	"github.com/NicholeGit/nade/framework/util/goroutine"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type GinServer struct {
	engine *gin.Engine
	server *http.Server
	isRun  bool
}

func (g *GinServer) Addr() string {
	return g.server.Addr
}

func (g *GinServer) Name() string {
	return "httpServer"
}

func NewGinServer(engine *gin.Engine, addr string) contract.IServer {
	ginServer := &GinServer{engine: engine, server: &http.Server{
		Handler:      engine,
		Addr:         addr,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}}
	return ginServer
}

func (g *GinServer) Start(_ context.Context) error {
	fmt.Println("GinServer Start")
	if g.server == nil {
		return errors.New("gin Server is not init")
	}
	g.isRun = true
	if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return errors.Wrap(err, "listen And Serve is error")
	}
	g.isRun = false
	return nil
}

func (g *GinServer) Stop(ctx context.Context) error {
	fmt.Println("GinServer Stop")
	err := goroutine.GraceStop(ctx, func() {
		_ = g.server.Shutdown(ctx)
	})
	if err != nil {
		return errors.Wrap(err, "GinServer stop is error")
	}
	return nil
}

func (g GinServer) WithAddr(addr string) error {
	if g.isRun {
		return errors.New("the address cannot be modified while the server is running")
	}
	g.server.Addr = addr
	return nil
}
