package bot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/service/session"
	"github.com/zhongxic/sellbot/internal/traceid"
	"github.com/zhongxic/sellbot/pkg/jieba"
)

type Service interface {
	Prologue(ctx context.Context, prologueDTO *PrologueDTO) (*InteractiveRespond, error)
}

type serviceImpl struct {
	options       Options
	testLoader    process.Loader
	releaseLoader process.Loader
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
	if s.options.ExtraDict == "" {
		tokenizer, err = jieba.NewDefaultTokenizer()
		slog.Debug(fmt.Sprintf("init default tokenizer cost [%v] ms", time.Since(start).Milliseconds()), traceId)
	} else {
		tokenizer, err = jieba.NewTokenizer(s.options.ExtraDict)
		slog.Debug(fmt.Sprintf("init tokenizer with extra dict [%v] cost [%v] ms",
			s.options.ExtraDict, time.Since(start).Milliseconds()), traceId)
	}
	return
}

type Options struct {
	ExtraDict string
}

func NewService(options Options, testLoader, releaseLoader process.Loader) Service {
	return &serviceImpl{
		options:       options,
		testLoader:    testLoader,
		releaseLoader: releaseLoader,
	}
}
