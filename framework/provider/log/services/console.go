package services

import (
	"os"

	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/contract"
)

// NadeConsoleLog 代表控制台输出
type NadeConsoleLog struct {
	NadeLog
}

// NewNadeConsoleLog 实例化HadeConsoleLog
func NewNadeConsoleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.IContainer)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	log := &NadeConsoleLog{}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	// 最重要的将内容输出到控制台
	log.SetOutput(os.Stdout)
	log.c = c
	return log, nil
}
