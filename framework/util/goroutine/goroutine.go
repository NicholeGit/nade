package goroutine

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"sync"

	"github.com/NicholeGit/nade/framework/util/errorx"
	"github.com/pkg/errors"

	"github.com/NicholeGit/nade/framework/command"
	"github.com/NicholeGit/nade/framework/contract"
)

// SafeGo 进行安全的goroutine调用
// 第一个参数是context接口，如果还实现了Container接口，且绑定了日志服务，则使用日志服务
// 第二个参数是匿名函数handler, 进行最终的业务逻辑
// SafeGo 函数并不会返回error，panic都会进入hade的日志服务
func SafeGo(ctx context.Context, handler func()) {
	var logger contract.ILog
	if c, ok := ctx.Value(command.ContextKey).(*command.ContextInfo); ok {
		container := c.Container()
		if container.IsBind(contract.LogKey) {
			logger = container.MustMake(contract.LogKey).(contract.ILog)
		}
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				buf := debug.Stack()
				buf = bytes.ReplaceAll(buf, []byte("\n"), []byte("\\n"))
				if logger != nil {
					logger.Error(ctx, "safe go handler panic", map[string]interface{}{
						"stack": string(buf),
						"err":   err,
					})
					fmt.Printf("panic\t%v\t%s", err, buf)
				} else {
					log.Printf("panic\t%v\t%s", err, buf)
				}
			}
		}()
		handler()
	}()
}

// SafeGoAndWait 进行并发安全并行调用
// 第一个参数是context接口，如果还实现了Container接口，且绑定了日志服务，则使用日志服务
// 第二个参数是匿名函数handlers数组, 进行最终的业务逻辑
// 返回handlers中任何一个错误（如果handlers中有业务逻辑返回错误）
func SafeGoAndWait(ctx context.Context, handlers ...func() error) error {
	var (
		wg     sync.WaitGroup
		once   sync.Once
		err    error
		logger contract.ILog
	)
	if c, ok := ctx.Value(command.ContextKey).(*command.ContextInfo); ok {
		container := c.Container()
		if container.IsBind(contract.LogKey) {
			logger = container.MustMake(contract.LogKey).(contract.ILog)
		}
	}

	for _, f := range handlers {
		wg.Add(1)
		go func(handler func() error) {
			defer func() {
				if err := recover(); err != nil {
					buf := debug.Stack()
					buf = bytes.ReplaceAll(buf, []byte("\n"), []byte("\\n"))
					if logger != nil {
						logger.Error(ctx, "panic", map[string]interface{}{
							"stack": string(buf),
							"err":   err,
						})
						fmt.Printf("panic\t%v\t%s", err, buf)
					} else {
						log.Printf("panic\t%v\t%s", err, buf)
					}
				}
				wg.Done()
			}()
			if e := handler(); e != nil {
				once.Do(func() {
					err = e
				})
			}
		}(f)
	}
	wg.Wait()
	return err
}

func GraceStop(ctx context.Context, stopFunc ...func() error) error {
	waitChan := make(chan int)
	var batchErr errorx.BatchError
	SafeGo(ctx, func() {
		defer close(waitChan)
		for _, v := range stopFunc {
			batchErr.Add(v())
		}
	})

	for {
		select {
		case <-waitChan:
			return batchErr.Err()
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "GraceStop is error")
		}
	}
}
