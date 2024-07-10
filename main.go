package main

import (
	"bwhNotify/config"
	"bwhNotify/logger"
	"bwhNotify/service"
	"bwhNotify/task"
	"bwhNotify/util"
	"flag"
	"os"
	"os/signal"
	"path"
	"syscall"
)

var configPath string

func init() {
	defaultPath := path.Join(util.GetCurrentAbPath(), "config.yml")
	flag.StringVar(&configPath, "c", defaultPath, "config file path")
	flag.Parse()
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	log := logger.NewLogger()

	// 读取配置文件
	if err := config.InitConf(configPath); err != nil {
		logger.Log().Error("init config file error", "err", err.Error())
		os.Exit(1)
	} else {
		logger.Log().Info("load config file success.", "config", config.Conf)
	}

	// 配置搬瓦工api服务
	if err := service.InitBwhApiService(config.Conf.BWHosts); err != nil {
		logger.Log().Error("init bwj api service error", "err", err.Error())
		os.Exit(2)
	}

	// 配置钉钉机器人服务
	if err := service.InitDingTalkService(config.Conf.DingTalk); err != nil {
		logger.Log().Error("init dingTalk service error", "err", err.Error())
		os.Exit(3)
	}

	// 启动定时任务
	t := task.NewTasker(log)
	if err := t.Start(); err != nil {
		logger.Log().Error("start cron task failed.", "err", err.Error())
		os.Exit(4)
	} else {
		logger.Log().Info("start cron task success.")
	}

	// 收到退出信号
	sig := <-sigs
	logger.Log().Info("Received signal.", "sig", sig.String())
	// 停止任务
	t.Stop()
	logger.Log().Info("Exiting...")
}
