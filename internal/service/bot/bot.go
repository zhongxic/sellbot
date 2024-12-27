package bot

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/service/session"
	"github.com/zhongxic/sellbot/internal/traceid"
	"github.com/zhongxic/sellbot/pkg/cache"
	"github.com/zhongxic/sellbot/pkg/jieba"
)

type Service interface {
	Prologue(ctx context.Context, prologueDTO *PrologueDTO) (*InteractiveRespond, error)
	Chat(ctx context.Context, chatDTO *ChatDTO) (*InteractiveRespond, error)
}

type serviceImpl struct {
	extraDict      string
	stopWords      string
	testLoader     process.Loader
	releaseLoader  process.Loader
	sessionCache   cache.Cache[string, *session.Session]
	tokenizerCache cache.Cache[string, *jieba.Tokenizer]
}

func (s *serviceImpl) Load(processId string, test bool) (*process.Process, error) {
	if test {
		return s.testLoader.Load(processId)
	}
	return s.releaseLoader.Load(processId)
}

func (s *serviceImpl) initSession(ctx context.Context, prologueDTO *PrologueDTO) *session.Session {
	sess := session.New()
	sess.ProcessId = prologueDTO.ProcessId
	sess.Variables = prologueDTO.Variables
	sess.Test = prologueDTO.Test
	slog.Debug(fmt.Sprintf("init session [%v]", sess.SessionId),
		"traceId", ctx.Value(traceid.TraceId{}), "prologueDTO", prologueDTO)
	return sess
}

func (s *serviceImpl) initTokenizer(ctx context.Context) (tokenizer *jieba.Tokenizer, err error) {
	start := time.Now()
	traceId := slog.Any("traceId", ctx.Value(traceid.TraceId{}))
	if s.extraDict == "" {
		tokenizer, err = jieba.NewDefaultTokenizer()
		slog.Debug(fmt.Sprintf("init default tokenizer cost [%v] ms", time.Since(start).Milliseconds()), traceId)
	} else {
		tokenizer, err = jieba.NewTokenizer(s.extraDict)
		slog.Debug(fmt.Sprintf("init tokenizer with extra dict [%v] cost [%v] ms",
			s.extraDict, time.Since(start).Milliseconds()), traceId)
	}
	return
}

func (s *serviceImpl) storeSession(sessionId string, sess *session.Session) {
	s.sessionCache.Set(sessionId, sess)
}

func (s *serviceImpl) storeTokenizer(sessionId string, tokenizer *jieba.Tokenizer) {
	s.tokenizerCache.Set(sessionId, tokenizer)
}

type Options struct {
	ExtraDict      string
	StopWords      string
	TestLoader     process.Loader
	ReleaseLoader  process.Loader
	SessionCache   cache.Cache[string, *session.Session]
	TokenizerCache cache.Cache[string, *jieba.Tokenizer]
}

func NewService(options Options) (Service, error) {
	if err := validate(options); err != nil {
		return nil, fmt.Errorf("create bot service failed: %w", err)
	}
	serve := &serviceImpl{
		extraDict:      options.ExtraDict,
		stopWords:      options.StopWords,
		testLoader:     options.TestLoader,
		releaseLoader:  options.ReleaseLoader,
		sessionCache:   options.SessionCache,
		tokenizerCache: options.TokenizerCache,
	}
	return serve, nil
}

func validate(options Options) error {
	if options.ExtraDict != "" {
		_, err := os.Stat(options.ExtraDict)
		if err != nil {
			return fmt.Errorf("extra dict [%v] not readable: %w", options.ExtraDict, err)
		}
	}
	if options.StopWords != "" {
		_, err := os.Stat(options.StopWords)
		if err != nil {
			return fmt.Errorf("stop words [%v] not readable: %w", options.StopWords, err)
		}
	}
	if options.TestLoader == nil {
		return fmt.Errorf("test loader is required")
	}
	if options.ReleaseLoader == nil {
		return fmt.Errorf("release loader is required")
	}
	if options.SessionCache == nil {
		return fmt.Errorf("session cache is required")
	}
	if options.TokenizerCache == nil {
		return fmt.Errorf("tokenizer cache is required")
	}
	return nil
}
