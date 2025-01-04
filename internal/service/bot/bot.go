package bot

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/zhongxic/sellbot/internal/service/bot/matcher"
	"github.com/zhongxic/sellbot/internal/service/bot/session"
	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/traceid"
	"github.com/zhongxic/sellbot/pkg/cache"
	"github.com/zhongxic/sellbot/pkg/jieba"
)

type Service interface {
	Prologue(ctx context.Context, prologueDTO *PrologueDTO) (*InteractiveRespond, error)
	Connect(ctx context.Context, connectDTO *SessionIdDTO) (*InteractiveRespond, error)
	Chat(ctx context.Context, chatDTO *ChatDTO) (*InteractiveRespond, error)
}

type serviceImpl struct {
	extraDict      string
	stopWords      []string
	processManager *process.Manager
	sessionManager session.Manager
	tokenizerCache cache.Cache[string, *jieba.Tokenizer]
	matcher        matcher.Matcher
}

func (s *serviceImpl) initSession(ctx context.Context, prologueDTO *PrologueDTO) *session.Session {
	sess := session.New()
	sess.ProcessId = prologueDTO.ProcessId
	sess.Variables = prologueDTO.Variables
	sess.Test = prologueDTO.Test
	slog.Info(fmt.Sprintf("init session [%v]", sess.Id),
		"traceId", ctx.Value(traceid.TraceId{}), "prologueDTO", prologueDTO)
	return sess
}

func (s *serviceImpl) initTokenizer(ctx context.Context) (tokenizer *jieba.Tokenizer, err error) {
	start := time.Now()
	traceId := slog.Any("traceId", ctx.Value(traceid.TraceId{}))
	if s.extraDict == "" {
		tokenizer, err = jieba.NewDefaultTokenizer()
		slog.Info(fmt.Sprintf("init default tokenizer cost [%v] ms", time.Since(start).Milliseconds()), traceId)
	} else {
		tokenizer, err = jieba.NewTokenizer(s.extraDict)
		slog.Info(fmt.Sprintf("init tokenizer with extra dict [%v] cost [%v] ms",
			s.extraDict, time.Since(start).Milliseconds()), traceId)
	}
	return
}

func (s *serviceImpl) storeSession(sessionId string, sess *session.Session) {
	s.sessionManager.Put(sessionId, sess)
}

func (s *serviceImpl) storeTokenizer(sessionId string, tokenizer *jieba.Tokenizer) {
	s.tokenizerCache.Set(sessionId, tokenizer)
}

func (s *serviceImpl) retrieveSession(sessionId string) (*session.Session, error) {
	sess := s.sessionManager.Get(sessionId)
	if sess == nil {
		return nil, fmt.Errorf("session [%s] not found", sessionId)
	}
	return sess, nil
}

func (s *serviceImpl) retrieveTokenizer(sessionId string) (*jieba.Tokenizer, error) {
	tokenizer, ok := s.tokenizerCache.Get(sessionId)
	if !ok {
		// TODO reload tokenizer
		return nil, fmt.Errorf("sessionId [%v]: tokenizer not found", sessionId)
	}
	return tokenizer, nil
}

type Options struct {
	ExtraDict      string
	StopWordsDict  string
	ProcessManager *process.Manager
	SessionManager session.Manager
	TokenizerCache cache.Cache[string, *jieba.Tokenizer]
	Matcher        matcher.Matcher
}

func NewService(options Options) (Service, error) {
	if err := validate(options); err != nil {
		return nil, fmt.Errorf("int bot service failed: %w", err)
	}
	stopWords, err := loadStopWords(options.StopWordsDict)
	if err != nil {
		return nil, fmt.Errorf("int bot service load stop words failed: %w", err)
	}
	serve := &serviceImpl{
		extraDict:      options.ExtraDict,
		stopWords:      stopWords,
		processManager: options.ProcessManager,
		sessionManager: options.SessionManager,
		tokenizerCache: options.TokenizerCache,
		matcher:        options.Matcher,
	}
	return serve, nil
}

func validate(options Options) error {
	if options.ExtraDict != "" {
		if _, err := os.Stat(options.ExtraDict); err != nil {
			return fmt.Errorf("stat extra dict [%v] failed: %w", options.ExtraDict, err)
		}
	}
	if options.StopWordsDict != "" {
		if _, err := os.Stat(options.StopWordsDict); err != nil {
			return fmt.Errorf("stat stop words [%v] failed: %w", options.StopWordsDict, err)
		}
	}
	if options.ProcessManager == nil {
		return fmt.Errorf("process manager is required")
	}
	if options.SessionManager == nil {
		return fmt.Errorf("session manager is required")
	}
	if options.TokenizerCache == nil {
		return fmt.Errorf("tokenizer cache is required")
	}
	if options.Matcher == nil {
		return fmt.Errorf("matcher is required")
	}
	return nil
}

func loadStopWords(stopWordsFile string) ([]string, error) {
	stopWords := make([]string, 0)
	if stopWordsFile != "" {
		f, err := os.Open(stopWordsFile)
		if err != nil {
			return stopWords, err
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				stopWords = append(stopWords, line)
			}
		}
	}
	return stopWords, nil
}
