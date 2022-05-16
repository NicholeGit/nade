package services

import (
	"os"
	"path/filepath"

	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/contract"
	"github.com/NicholeGit/nade/framework/util"
	"github.com/pkg/errors"
)

type NadeSingleLog struct {
	NadeLog

	folder string
	file   string
	fd     *os.File
}

// NewHadeSingleLog params sequence: level, ctxFielder, Formatter, map[string]interface(folder/file)
func NewHadeSingleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.IContainer)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	appService := c.MustMake(contract.AppKey).(contract.IApp)
	configService := c.MustMake(contract.ConfigKey).(contract.IConfig)

	log := &NadeSingleLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	folder := appService.LogFolder()
	if configService.IsExist("log.folder") {
		folder = configService.GetString("log.folder")
	}
	log.folder = folder
	if !util.Exists(folder) {
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return nil, errors.Wrap(err, "MkdirAll error")
		}
	}

	log.file = "nade.log"
	if configService.IsExist("log.file") {
		log.file = configService.GetString("log.file")
	}
	fd, err := os.OpenFile(filepath.Join(log.folder, log.file), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "open log file err")
	}
	log.fd = fd

	log.SetOutput(fd)
	log.c = c
	log.fd = fd

	return log, nil
}
