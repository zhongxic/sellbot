package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/zhongxic/sellbot/config"
	"github.com/zhongxic/sellbot/internal/routes"
	"github.com/zhongxic/sellbot/pkg/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config/config.yaml", "the config file in yaml format")
	flag.Parse()
}

func main() {
	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}
	if err := logger.Init(cfg.Logging); err != nil {
		log.Fatal(err)
	}
	r := routes.Init()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Server.Port),
		Handler: r,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		logger.Info("server is shutting down")
		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Info("server shutdown", "error", err)
		}
		close(idleConnsClosed)
	}()

	logger.Info("server started", "config", cfg)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server listen and serve: %v", err)
	}

	<-idleConnsClosed
}
