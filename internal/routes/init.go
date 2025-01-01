package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhongxic/sellbot/config"
	botctl "github.com/zhongxic/sellbot/internal/controller/bot"
	"github.com/zhongxic/sellbot/internal/controller/ping"
	botserve "github.com/zhongxic/sellbot/internal/service/bot"
	"github.com/zhongxic/sellbot/internal/service/bot/matcher"
	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/service/session"
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
	testProcessLoader := process.NewCachedLoader(process.NewFileLoader(cfg.Process.Directory.Test), testProcessCache)
	releaseProcessLoader := process.NewCachedLoader(process.NewFileLoader(cfg.Process.Directory.Release), releaseProcessCache)
	processManager := &process.Manager{TestProcessLoader: releaseProcessLoader, ReleaseProcessLoader: testProcessLoader}
	sessionCache := cache.NewCache[*session.Session](cache.Options{
		DefaultExpiration: time.Duration(cfg.Session.Cache.Expiration) * time.Second,
		CleanupInterval:   time.Duration(cfg.Session.Cache.CleanupInterval) * time.Second,
	})
	tokenizerCache := cache.NewCache[*jieba.Tokenizer](cache.Options{
		DefaultExpiration: time.Duration(cfg.Session.Cache.Expiration) * time.Second,
		CleanupInterval:   time.Duration(cfg.Session.Cache.CleanupInterval) * time.Second,
	})
	botOptions := botserve.Options{
		ExtraDict:      cfg.Tokenizer.ExtraDict,
		StopWords:      cfg.Tokenizer.StopWords,
		ProcessManager: processManager,
		SessionCache:   sessionCache,
		TokenizerCache: tokenizerCache,
		Matcher:        matcher.DefaultChainedMatcher,
	}
	return botserve.NewService(botOptions)
}
