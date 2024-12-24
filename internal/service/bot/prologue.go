package bot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/zhongxic/sellbot/internal/service/bot/matcher"
	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/service/process/helper"
	"github.com/zhongxic/sellbot/internal/service/session"
	"github.com/zhongxic/sellbot/internal/traceid"
	"github.com/zhongxic/sellbot/pkg/jieba"
)

func (s *serviceImpl) Prologue(ctx context.Context, prologueDTO *PrologueDTO) (*InteractiveRespond, error) {
	slog.Info("start process prologue", "traceId", ctx.Value(traceid.TraceId{}))
	loadedProcess, err := s.Load(prologueDTO.ProcessId, prologueDTO.Test)
	if err != nil {
		return nil, err
	}
	if err := loadedProcess.Validate(); err != nil {
		return nil, err
	}
	if err := validateVariables(prologueDTO.Variables, loadedProcess.Variables); err != nil {
		return nil, err
	}
	currentSession := s.initSession(ctx, prologueDTO)
	tokenizer, err := s.initTokenizer(ctx)
	if err != nil {
		return nil, err
	}
	processHelper := helper.New(loadedProcess)
	startDomain, err := processHelper.FindStartDomain()
	if err != nil {
		return nil, err
	}
	loadUserDict(tokenizer, processHelper)
	matchContext := matcher.NewContext(currentSession, loadedProcess)
	matchContext.AddMatchedPath(matcher.MatchedPath{Domain: startDomain.Name, Branch: process.BranchNameEnter})
	answerDTO, err := makeAnswer(ctx, matchContext)
	if err != nil {
		return nil, err
	}
	statPaths := convertMatchedPathsToStatPaths(matchContext.MatchedPaths)
	currentSession.UpdateStat(statPaths)
	intentionRules := []process.IntentionRule{processHelper.GetDefaultIntentionRule()}
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

func loadUserDict(tokenizer *jieba.Tokenizer, processHelper *helper.Helper) error {
	globalKeywords := processHelper.GetGlobalKeywords()
	for _, keyword := range globalKeywords {
		tokenizer.AddWord(keyword, 1)
	}
	startDomain, err := processHelper.FindStartDomain()
	if err != nil {
		return nil
	}
	startDomainKeywords := processHelper.GetDomainKeywords(startDomain.Name)
	for _, keyword := range startDomainKeywords {
		tokenizer.AddWord(keyword, 1)
	}
	intentionKeyword := processHelper.GetIntentionKeywords()
	for _, keyword := range intentionKeyword {
		tokenizer.AddWord(keyword, 1)
	}
	return nil
}

func convertMatchedPathsToStatPaths(matchedPaths []matcher.MatchedPath) []session.StatPath {
	// TODO impl-me convertMatchedPathToStatPath
	return make([]session.StatPath, 0)
}
