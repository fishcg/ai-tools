package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fish/ai-tools/config"
	"github.com/fish/ai-tools/controllers"
	"github.com/fish/ai-tools/logger"
	"github.com/fish/ai-tools/server"
	"github.com/fish/ai-tools/service"
)

var configPath = "./config.yml"

// 初始化操作
func init() {
	configFile := flag.String("c", configPath, "config file path")
	flag.Parse()
	err := config.LoadFromYML(*configFile)
	if err != nil {
		logger.Fatalf("load config error: %v", err)
	}

	err = service.Init(config.Conf)
	if err != nil {
		logger.Infof("init service error: %v", err)
	}
}

func main() {
	srv := server.NewHTTPServer(&config.Conf.HTTP)
	srv.Run(controllers.MountGame)
	shutdownServer(srv)
}

func shutdownServer(srv server.Server) {
	quit := make(chan os.Signal, 1)
	// kill send SIGTERM; kill -2 send SIGINT; kill -9 can't be catch
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Servers ...")
	// NOTICE: 超时时间应该基本与最长请求时间时间相等
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server ", "Shutdown: ", err)
	}

	logger.Info("All Servers exiting")
}
