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
	loadedProcess, err := s.processManager.Load(currentSession.ProcessId, currentSession.Test)
	if err != nil {
		return nil, fmt.Errorf("load process failed: %w", err)
	}
	matchContext := matcher.NewContext(currentSession, loadedProcess)
	matchContext.Sentence = chatDTO.Sentence
	matchContext.Segments = cutAll(ctx, tokenizer, s.stopWords, chatDTO.Sentence)
	matchContext.Silence = chatDTO.Silence
	matchContext.Interruption = chatDTO.Interruption
	_, err = s.matcher.Match(ctx, matchContext)
	if err != nil {
		slog.Error(fmt.Sprintf("sessionId [%v]: process matching failed: %v", currentSession.Id, err),
			slog.Any("traceId", traceid.TraceId{}))
		processHelper := process.NewHelper(loadedProcess)
		domain, err := processHelper.GetCommonDialog(process.DomainTypeDialogEndException)
		if err != nil {
			return nil, fmt.Errorf("get common dialog [%v] failed: %w", process.DomainTypeDialogEndException, err)
		}
		matchedPath := matcher.MatchedPath{Domain: domain.Name, Branch: process.BranchNameEnter}
		slog.Info(fmt.Sprintf("jump to domain [%v] branch [%v] due to match error occurred",
			matchedPath.Domain, matchedPath.Branch),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		matchContext.AddMatchedPath(matchedPath)
	}
	answerDTO, err := makeAnswer(ctx, matchContext)
	if err != nil {
		return nil, fmt.Errorf("make answer failed: %w", err)
	}
	previousMainProcessDomain := currentSession.CurrentMainProcessDomain
	matchContext.UpdateSessionStat()
	intentionRules := analyzeIntention(currentSession, loadedProcess)
	reloadKeywords(tokenizer, loadedProcess, previousMainProcessDomain, currentSession.CurrentMainProcessDomain)
	s.storeSession(currentSession.Id, currentSession)
	s.storeTokenizer(currentSession.Id, tokenizer)
	return makeRespond(matchContext, answerDTO, intentionRules), nil
}

func cutAll(ctx context.Context, tokenizer *jieba.Tokenizer, stopWords []string, sentence string) []string {
	cuts := make([]string, 0)
	words := tokenizer.CutAll(sentence)
	slog.Info(fmt.Sprintf("cut sentence [%v] in to words [%v]", sentence, words), "traceId", ctx.Value(traceid.TraceId{}))
	for _, word := range words {
		// TODO remove stop words
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

func reloadKeywords(tokenizer *jieba.Tokenizer, loadedProcess *process.Process,
	previousMainProcessDomain, lastMainProcessDomain string) {
	// TODO reload keywords
}
