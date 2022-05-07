package services

import (
	"context"
	"io"
	golog "log"
	"time"

	"github.com/pkg/errors"

	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/contract"
	"github.com/NicholeGit/nade/framework/provider/log/formatter"
)

// NadeLog 的通用实例
type NadeLog struct {
	// 五个必要参数
	level      contract.LogLevel    // 日志级别
	formatter  contract.Formatter   // 日志格式化方法
	ctxFielder contract.CtxFielder  // ctx获取上下文字段
	output     io.Writer            // 输出
	c          framework.IContainer // 容器
}

// IsLevelEnable 判断这个级别是否可以打印
func (log *NadeLog) IsLevelEnable(level contract.LogLevel) bool {
	return level <= log.level
}

// logf 为打印日志的核心函数
func (log *NadeLog) logf(level contract.LogLevel, ctx context.Context, msg string, fields map[string]interface{}) error {
	// 先判断日志级别
	if !log.IsLevelEnable(level) {
		return nil
	}

	// 使用ctxFielder 获取context中的信息
	fs := fields
	if log.ctxFielder != nil {
		t := log.ctxFielder(ctx)
		if t != nil {
			for k, v := range t {
				fs[k] = v
			}
		}
	}

	// 如果绑定了trace服务，获取trace信息
	// if log.c != nil && log.c.IsBind(contract.TraceKey) {
	// 	tracer := log.c.MustMake(contract.TraceKey).(contract.Trace)
	// 	tc := tracer.GetTrace(ctx)
	// 	if tc != nil {
	// 		maps := tracer.ToMap(tc)
	// 		for k, v := range maps {
	// 			fs[k] = v
	// 		}
	// 	}
	// }

	// 将日志信息按照formatter序列化为字符串
	if log.formatter == nil {
		log.formatter = formatter.TextFormatter
	}
	ct, err := log.formatter(level, time.Now(), msg, fs)
	if err != nil {
		return err
	}

	// 如果是panic级别，则使用log进行panic
	if level == contract.PanicLevel {
		golog.Panicln(string(ct))
		return nil
	}

	// 通过output进行输出
	_, err = log.output.Write(ct)
	if err != nil {
		return errors.Wrap(err, "Write error")
	}
	_, err = log.output.Write([]byte("\r\n"))
	if err != nil {
		return errors.Wrap(err, "Write error")
	}
	return nil
}

// SetOutput 设置output
func (log *NadeLog) SetOutput(output io.Writer) {
	log.output = output
}

// Panic 输出panic的日志信息
func (log *NadeLog) Panic(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.PanicLevel, ctx, msg, fields)
}

// Fatal will add fatal record which contains msg and fields
func (log *NadeLog) Fatal(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.FatalLevel, ctx, msg, fields)
}

// Error will add error record which contains msg and fields
func (log *NadeLog) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.ErrorLevel, ctx, msg, fields)
}

// Warn will add warn record which contains msg and fields
func (log *NadeLog) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.WarnLevel, ctx, msg, fields)
}

// Info 会打印出普通的日志信息
func (log *NadeLog) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.InfoLevel, ctx, msg, fields)
}

// Debug will add debug record which contains msg and fields
func (log *NadeLog) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.DebugLevel, ctx, msg, fields)
}

// Trace will add trace info which contains msg and fields
func (log *NadeLog) Trace(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.TraceLevel, ctx, msg, fields)
}

// SetLevel set log level, and higher level will be recorded
func (log *NadeLog) SetLevel(level contract.LogLevel) {
	log.level = level
}

// SetCtxFielder will get fields from context
func (log *NadeLog) SetCtxFielder(handler contract.CtxFielder) {
	log.ctxFielder = handler
}

// SetFormatter will set formatter handler will covert data to string for recording
func (log *NadeLog) SetFormatter(formatter contract.Formatter) {
	log.formatter = formatter
}
