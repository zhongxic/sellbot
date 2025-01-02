package bot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/zhongxic/sellbot/internal/service/bot/matcher"
	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/traceid"
	"github.com/zhongxic/sellbot/pkg/jieba"
)

func (s *serviceImpl) Prologue(ctx context.Context, prologueDTO *PrologueDTO) (*InteractiveRespond, error) {
	slog.Info("start process prologue", "traceId", ctx.Value(traceid.TraceId{}))
	loadedProcess, err := s.processManager.Load(prologueDTO.ProcessId, prologueDTO.Test)
	if err != nil {
		return nil, fmt.Errorf("load process failed: %w", err)
	}
	if err := loadedProcess.Validate(); err != nil {
		return nil, fmt.Errorf("process validate failed: %w", err)
	}
	if err := validateVariables(prologueDTO.Variables, loadedProcess.Variables); err != nil {
		return nil, fmt.Errorf("variables validate failed: %w ", err)
	}
	currentSession := s.initSession(ctx, prologueDTO)
	tokenizer, err := s.initTokenizer(ctx)
	if err != nil {
		return nil, fmt.Errorf("init tokenizer failed: %w", err)
	}
	processHelper := process.NewHelper(loadedProcess)
	if err := loadUserDict(tokenizer, processHelper); err != nil {
		return nil, fmt.Errorf("load user dict failed: %w", err)
	}
	startDomain, err := processHelper.GetStartDomain()
	if err != nil {
		return nil, fmt.Errorf("get start domain failed: %w", err)
	}
	matchContext := matcher.NewContext(currentSession, loadedProcess)
	matchContext.AddMatchedPath(matcher.MatchedPath{Domain: startDomain.Name, Branch: process.BranchNameEnter})
	answerDTO, err := makeAnswer(ctx, matchContext)
	if err != nil {
		return nil, fmt.Errorf("make answer failed: %w", err)
	}
	if err := matchContext.UpdateSessionStat(); err != nil {
		return nil, fmt.Errorf("update session stat failed: %w", err)
	}
	intentionRules := []process.IntentionRule{processHelper.GetDefaultIntentionRule()}
	s.storeSession(currentSession.Id, currentSession)
	s.storeTokenizer(currentSession.Id, tokenizer)
	return makeRespond(matchContext, answerDTO, intentionRules), nil
}

func validateVariables(actual map[string]string, expected []process.Variable) error {
	params := actual
	if actual == nil {
		params = make(map[string]string)
	}
	variables := expected
	if expected == nil {
		variables = make([]process.Variable, 0)
	}
	if len(params) != len(variables) {
		return fmt.Errorf("process variables not matched expected [%d] actual [%d]", len(variables), len(params))
	}
	messages := make([]string, 0)
	for _, variable := range variables {
		if _, ok := params[variable.Code]; !ok {
			message := fmt.Sprintf("process variable [%s] is required", variable.Code)
			messages = append(messages, message)
		}
	}
	if len(messages) > 0 {
		return errors.New(strings.Join(messages, ", "))
	}
	return nil
}

func loadUserDict(tokenizer *jieba.Tokenizer, processHelper *process.Helper) error {
	globalKeywords := processHelper.GetGlobalKeywords()
	for _, keyword := range globalKeywords {
		tokenizer.AddWord(keyword, 1)
	}
	startDomain, err := processHelper.GetStartDomain()
	if err != nil {
		return fmt.Errorf("get start domain failed: %w", err)
	}
	startDomainKeywords := processHelper.GetDomainKeywords(startDomain.Name)
	for _, keyword := range startDomainKeywords {
		tokenizer.AddWord(keyword, 1)
	}
	return nil
}
