package main

import (
	"flag"

	"github.com/wwq1988/leaf/pkg/conf"
	"github.com/wwq1988/leaf/pkg/logger"

	"github.com/wwq1988/leaf/pkg/generator"
	"github.com/wwq1988/leaf/pkg/server"
	"github.com/wwq1988/leaf/pkg/service"
	"github.com/wwq1988/leaf/pkg/util"
	"go.uber.org/zap"
)

var (
	cfgPath = flag.String("conf", "./conf.toml", "-conf=./conf.toml")
)

func main() {
	flag.Parse()
	if *cfgPath == "" {
		flag.Usage()
		return
	}
	if err := conf.Parse(*cfgPath); err != nil {
		logger.Fatal("parse conf", zap.String("conf", *cfgPath), zap.Any("err", err))
	}

	logger := logger.New()

	generator, err := generator.New()
	if err != nil {
		logger.Fatal("new generator", zap.Any("err", err))
	}

	service := service.New(generator, logger)
	srv := server.New(service)
	if err := srv.Init(); err != nil {
		logger.Fatal("init server", zap.Any("err", err))
	}
	util.OnExit(srv.Stop, func() { logger.Sync() })
	if err := srv.Start(); err != nil {
		logger.Fatal("server start", zap.Any("err", err))
	}
}
