package bot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/service/session"
	"github.com/zhongxic/sellbot/internal/traceid"
	"github.com/zhongxic/sellbot/pkg/cache"
	"github.com/zhongxic/sellbot/pkg/jieba"
)

type Service interface {
	Prologue(ctx context.Context, prologueDTO *PrologueDTO) (*InteractiveRespond, error)
}

type serviceImpl struct {
	extraDict      string
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

type Options struct {
	ExtraDict      string
	TestLoader     process.Loader
	ReleaseLoader  process.Loader
	SessionCache   cache.Cache[string, *session.Session]
	TokenizerCache cache.Cache[string, *jieba.Tokenizer]
}

func NewService(options Options) Service {
	return &serviceImpl{
		extraDict:      options.ExtraDict,
		testLoader:     options.TestLoader,
		releaseLoader:  options.ReleaseLoader,
		sessionCache:   options.SessionCache,
		tokenizerCache: options.TokenizerCache,
	}
}
