package services

import (
	"io"

	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/contract"
)

type NadeCustomLog struct {
	NadeLog
}

// NewNadeCustomLog output 需要自己定义
func NewNadeCustomLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.IContainer)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)
	output := params[4].(io.Writer)

	log := &NadeConsoleLog{}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	log.SetOutput(output)
	log.c = c
	return log, nil
}
