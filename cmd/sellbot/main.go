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
	"syscall"

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
		return
	}

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}
	if err := logger.Init(cfg.Logging); err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = logger.Close()
	}()
	r, err := routes.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Server.Port),
		Handler: r,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT)
		<-sigint

		log.Println("server shutting down")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("server shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("server started with config: %+v", cfg)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server listen and serve: %v", err)
	}

	<-idleConnsClosed
}
