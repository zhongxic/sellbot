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
	stopWords      []string
	processManager *process.Manager
	sessionCache   cache.Cache[string, *session.Session]
	tokenizerCache cache.Cache[string, *jieba.Tokenizer]
	matcher        matcher.Matcher
}

func (s *serviceImpl) initSession(ctx context.Context, prologueDTO *PrologueDTO) *session.Session {
	sess := session.New()
	sess.ProcessId = prologueDTO.ProcessId
	sess.Variables = prologueDTO.Variables
	sess.Test = prologueDTO.Test
	slog.Info(fmt.Sprintf("init session [%v]", sess.SessionId),
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
	s.sessionCache.Set(sessionId, sess)
}

func (s *serviceImpl) storeTokenizer(sessionId string, tokenizer *jieba.Tokenizer) {
	s.tokenizerCache.Set(sessionId, tokenizer)
}

func (s *serviceImpl) retrieveSession(sessionId string) (*session.Session, error) {
	sess, ok := s.sessionCache.Get(sessionId)
	if !ok {
		return nil, fmt.Errorf("sessionId [%v]: session not found", sessionId)
	}
	return sess, nil
}

func (s *serviceImpl) retrieveTokenizer(sessionId string) (*jieba.Tokenizer, error) {
	tokenizer, ok := s.tokenizerCache.Get(sessionId)
	if !ok {
		return nil, fmt.Errorf("sessionId [%v]: tokenizer not found", sessionId)
	}
	return tokenizer, nil
}

type Options struct {
	ExtraDict      string
	StopWords      string
	ProcessManager *process.Manager
	SessionCache   cache.Cache[string, *session.Session]
	TokenizerCache cache.Cache[string, *jieba.Tokenizer]
	Matcher        matcher.Matcher
}

func NewService(options Options) (Service, error) {
	if err := validate(options); err != nil {
		return nil, fmt.Errorf("create bot service failed: %w", err)
	}
	serve := &serviceImpl{
		extraDict:      options.ExtraDict,
		processManager: options.ProcessManager,
		sessionCache:   options.SessionCache,
		tokenizerCache: options.TokenizerCache,
		matcher:        options.Matcher,
	}
	stopWords, err := loadStopWords(options.StopWords)
	if err != nil {
		return nil, fmt.Errorf("create bot service load stop words failed: %w", err)
	}
	serve.stopWords = stopWords
	return serve, nil
}

func loadStopWords(stopWordsFile string) ([]string, error) {
	stopWords := make([]string, 0)
	if stopWordsFile != "" {
		f, err := os.Open(stopWordsFile)
		if err != nil {
			return stopWords, fmt.Errorf("load stop words [%v] failed: %w", stopWordsFile, err)
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
	if options.ProcessManager == nil {
		return fmt.Errorf("process manager is required")
	}
	if options.SessionCache == nil {
		return fmt.Errorf("session cache is required")
	}
	if options.TokenizerCache == nil {
		return fmt.Errorf("tokenizer cache is required")
	}
	if options.Matcher == nil {
		return fmt.Errorf("matcher is required")
	}
	return nil
}
