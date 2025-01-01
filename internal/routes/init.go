package routes

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/config"
	botctl "github.com/zhongxic/sellbot/internal/controller/bot"
	"github.com/zhongxic/sellbot/internal/controller/ping"
	botserve "github.com/zhongxic/sellbot/internal/service/bot"
	"github.com/zhongxic/sellbot/internal/service/bot/matcher"
	"github.com/zhongxic/sellbot/internal/service/bot/session"
	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/pkg/cache"
	"github.com/zhongxic/sellbot/pkg/jieba"
	"github.com/zhongxic/sellbot/pkg/middleware"
)

func Init(cfg *config.Config) (*gin.Engine, error) {
	return initRoutes(cfg)
}

func initRoutes(cfg *config.Config) (*gin.Engine, error) {
	r := gin.New()
	registerMiddleware(r)
	err := registerRoutes(r, cfg)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func registerMiddleware(r *gin.Engine) {
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
}

func registerRoutes(r *gin.Engine, cfg *config.Config) error {
	pingController := ping.NewController()
	botService, err := newBotService(cfg)
	if err != nil {
		return err
	}
	botController := botctl.NewController(botService)
	r.GET("/ping", pingController.Ping)
	r.POST("/prologue", botController.Prologue)
	r.POST("/chat", botController.Chat)
	return nil
}

func newBotService(cfg *config.Config) (botserve.Service, error) {
	testProcessCache := cache.NewCache[*process.Process](cache.Options{
		DefaultExpiration: time.Duration(cfg.Process.Cache.Expiration) * time.Second,
		CleanupInterval:   time.Duration(cfg.Process.Cache.CleanupInterval) * time.Second,
	})
	releaseProcessCache := cache.NewCache[*process.Process](cache.Options{
		DefaultExpiration: time.Duration(cfg.Process.Cache.Expiration) * time.Second,
		CleanupInterval:   time.Duration(cfg.Process.Cache.CleanupInterval) * time.Second,
	})
	processManager := &process.Manager{
		TestProcessLoader:    process.NewCachedLoader(process.NewFileLoader(cfg.Process.Directory.Test), testProcessCache),
		ReleaseProcessLoader: process.NewCachedLoader(process.NewFileLoader(cfg.Process.Directory.Release), releaseProcessCache),
	}
	sessionManager, err := session.NewManager(session.Options{
		Repository: cfg.Session.Repository,
		Expiration: time.Duration(cfg.Session.Expiration) * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("create session manager failed: %w", err)
	}
	tokenizerCache := cache.NewCache[*jieba.Tokenizer](cache.Options{
		DefaultExpiration: time.Duration(cfg.Session.Expiration) * time.Second,
		CleanupInterval:   session.DefaultCleanupInterval,
	})
	botOptions := botserve.Options{
		ExtraDict:      cfg.Tokenizer.ExtraDict,
		StopWordsDict:  cfg.Tokenizer.StopWordsDict,
		ProcessManager: processManager,
		SessionManager: sessionManager,
		TokenizerCache: tokenizerCache,
		Matcher:        matcher.DefaultChainedMatcher,
	}
	return botserve.NewService(botOptions)
}
