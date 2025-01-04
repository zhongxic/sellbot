package bot

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"slices"

	"github.com/zhongxic/sellbot/internal/service/bot/matcher"
	"github.com/zhongxic/sellbot/internal/service/bot/session"
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
	segments := cutAll(ctx, tokenizer, s.stopWords, chatDTO.Sentence)
	matchContext := matcher.NewContext(currentSession, loadedProcess)
	matchContext.Sentence = chatDTO.Sentence
	matchContext.Segments = segments
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
	if err := matchContext.UpdateSessionStat(); err != nil {
		return nil, fmt.Errorf("update session stat failed: %w", err)
	}
	env := assembleIntentionAnalyzeEnv(chatDTO.Sentence, segments, currentSession)
	intentionRules := analyzeIntention(ctx, env, loadedProcess.Intentions.IntentionRules)
	reloadKeywords(tokenizer, loadedProcess, currentSession)
	s.storeSession(currentSession.Id, currentSession)
	s.storeTokenizer(currentSession.Id, tokenizer)
	return makeRespond(matchContext, answerDTO, intentionRules), nil
}

func cutAll(ctx context.Context, tokenizer *jieba.Tokenizer, stopWords []string, sentence string) []string {
	cuts := make([]string, 0)
	words := tokenizer.CutAll(sentence)
	slog.Info(fmt.Sprintf("cut sentence [%v] in to words [%v]", sentence, words), "traceId", ctx.Value(traceid.TraceId{}))
	for _, word := range words {
		if !slices.Contains(stopWords, word) && reHan.MatchString(word) {
			cuts = append(cuts, word)
		}
	}
	slog.Info(fmt.Sprintf("cut sentence [%v] in to words [%v] after stop words removed", sentence, cuts),
		"traceId", ctx.Value(traceid.TraceId{}))
	return cuts
}

func reloadKeywords(tokenizer *jieba.Tokenizer, loadedProcess *process.Process, currentSession *session.Session) {
	if currentSession.PreviousMainProcessDomain == currentSession.CurrentMainProcessDomain {
		return
	}
	processHelper := process.NewHelper(loadedProcess)
	previousDomainKeywords := processHelper.GetDomainKeywords(currentSession.PreviousMainProcessDomain)
	for _, keyword := range previousDomainKeywords {
		tokenizer.AddWord(keyword, -1)
	}
	currentDomainKeywords := processHelper.GetDomainKeywords(currentSession.CurrentMainProcessDomain)
	for _, keyword := range currentDomainKeywords {
		tokenizer.AddWord(keyword, 1)
	}
}
