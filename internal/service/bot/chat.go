package bot

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"slices"

	"github.com/zhongxic/sellbot/internal/service/bot/matcher"
	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/traceid"
	"github.com/zhongxic/sellbot/pkg/jieba"
)

var reHan = regexp.MustCompile(`\p{Han}+`)

func (s *serviceImpl) Chat(ctx context.Context, chatDTO *ChatDTO) (*InteractiveRespond, error) {
	slog.Info("start process chat", "traceId", ctx.Value(traceid.TraceId{}))
	currentSession, err := s.retrieveSession(chatDTO.SessionId)
	if err != nil {
		return nil, fmt.Errorf("retrieve session failed: %w", err)
	}
	tokenizer, err := s.retrieveTokenizer(chatDTO.SessionId)
	if err != nil {
		return nil, fmt.Errorf("retrieve tokenizer failed: %w", err)
	}
	loadedProcess, err := s.loadProcess(currentSession.ProcessId, currentSession.Test)
	if err != nil {
		return nil, fmt.Errorf("load process failed: %w", err)
	}
	matchContext := matcher.NewContext(currentSession, loadedProcess)
	matchContext.Sentence = chatDTO.Sentence
	matchContext.Segments = cutAll(ctx, tokenizer, s.stopWords, chatDTO.Sentence)
	matchContext.Silence = chatDTO.Silence
	matchContext.Interruption = chatDTO.Interruption
	// TODO match in process
	answerDTO, err := makeAnswer(ctx, matchContext)
	if err != nil {
		return nil, fmt.Errorf("make answer failed: %w", err)
	}
	statPaths := convertMatchedPathsToStatPaths(matchContext.MatchedPaths)
	currentSession.UpdateStat(statPaths)
	// TODO intention analysis
	intentionRules := make([]process.IntentionRule, 0)
	s.storeSession(currentSession.SessionId, currentSession)
	s.storeTokenizer(currentSession.SessionId, tokenizer)
	return makeRespond(matchContext, answerDTO, intentionRules), nil
}

func cutAll(ctx context.Context, tokenizer *jieba.Tokenizer, stopWords []string, sentence string) []string {
	cuts := make([]string, 0)
	words := tokenizer.CutAll(sentence)
	slog.Info(fmt.Sprintf("cut sentence [%v] in to words [%v]", sentence, words), "traceId", ctx.Value(traceid.TraceId{}))
	for _, word := range words {
		contains := slices.Index(stopWords, word) == -1
		if contains {
			continue
		}
		if !reHan.MatchString(word) {
			continue
		}
		cuts = append(cuts, word)
	}
	slog.Info(fmt.Sprintf("cut sentence [%v] in to words [%v] after stop words removed", sentence, cuts),
		"traceId", ctx.Value(traceid.TraceId{}))
	return cuts
}
