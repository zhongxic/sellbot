package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/zhongxic/sellbot/config"
	"github.com/zhongxic/sellbot/internal/routes"
	"github.com/zhongxic/sellbot/pkg/logger"
)

const version = "1.0-SNAPSHOT"

var (
	showVersion bool
	configFile  string
)

func init() {
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.StringVar(&configFile, "config", "config/config.yaml", "config file in yaml format")
	flag.Parse()
}

func main() {
	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}
	if err := logger.Init(cfg.Logging); err != nil {
		log.Fatal(err)
	}
	defer logger.Close()
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

		slog.Info("server is shutting down")
		if err := srv.Shutdown(context.Background()); err != nil {
			slog.Info("server shutdown", slog.Any("error", err))
		}
		close(idleConnsClosed)
	}()

	slog.Info("server started", slog.Any("config", cfg))

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server listen and serve: %v", err)
	}

	<-idleConnsClosed
}
