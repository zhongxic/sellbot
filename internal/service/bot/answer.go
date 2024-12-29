package bot

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/zhongxic/sellbot/internal/service/bot/matcher"
	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/service/process/helper"
	"github.com/zhongxic/sellbot/internal/traceid"
)

func makeAnswer(ctx context.Context, matchContext *matcher.Context) (AnswerDTO, error) {
	if matchContext == nil {
		return AnswerDTO{}, fmt.Errorf("make answer failed due to nil match context")
	}
	matchedPath, err := matchContext.GetLastMatchedPath()
	if err != nil {
		return AnswerDTO{}, fmt.Errorf("get last matched path failed: %w", err)
	}
	traceId := slog.Any("traceId", ctx.Value(traceid.TraceId{}))
	slog.Info(fmt.Sprintf("sessionId [%v]: matched domain [%v] branch [%v]",
		matchContext.Session.SessionId, matchedPath.Domain, matchedPath.Branch), traceId)
	processHelper := helper.New(matchContext.Process)
	domain, err := processHelper.GetDomain(matchedPath.Domain)
	if err != nil {
		return AnswerDTO{}, fmt.Errorf("get domain failed: %w", err)
	}
	branch, err := processHelper.GetBranch(matchedPath.Domain, matchedPath.Branch)
	if err != nil {
		return AnswerDTO{}, fmt.Errorf("get branch failed: %w", err)
	}
	hitCount := matchContext.Session.GetDomainBranchHitCount(matchedPath.Domain, matchedPath.Branch)
	isExceed := hitCount >= len(branch.Responses) && domain.Category != process.DomainCategoryMainProcess
	slog.Info(fmt.Sprintf("sessionId [%v]: domain [%v] branch [%v] hitCount [%v] isExceed [%v]",
		matchContext.Session.SessionId, matchedPath.Domain, matchedPath.Branch, hitCount, isExceed), traceId)
	if isExceed {
		nextDomain := ""
		if branch.EnableExceedJump && branch.Next != "" {
			nextDomain = branch.Next
		}
		slog.Info(fmt.Sprintf("sessionId [%v]: jump to domain [%v] due to hitCount exceed",
			matchContext.Session.SessionId, matchedPath.Branch), traceId)
		return autoJump(ctx, matchContext, nextDomain)
	}
	response := branch.Responses[hitCount%len(branch.Responses)]
	if response.EnableAutoJump && response.Next != "" {
		slog.Info(fmt.Sprintf("sessionId [%v]: jump to domain [%v] due to domain [%v] branch [%v] auto jump enabled",
			matchContext.Session.SessionId, response.Next, matchedPath.Domain, matchedPath.Branch), traceId)
		return autoJump(ctx, matchContext, response.Next)
	}
	return AnswerDTO{Text: response.Text, Audio: response.Audio}, nil
}

func autoJump(ctx context.Context, matchContext *matcher.Context, nextDomain string) (AnswerDTO, error) {
	matchedPath := matcher.MatchedPath{Domain: nextDomain, Branch: process.BranchNameEnter}
	if matchedPath.Domain == "" {
		processHelper := helper.New(matchContext.Process)
		endFailDomain, err := processHelper.GetCommonDialogDomain(process.DomainTypeDialogEndFail)
		if err != nil {
			return AnswerDTO{}, fmt.Errorf("find common dialog domain [%v] failed: %w", process.DomainTypeDialogEndFail, err)
		}
		matchedPath.Domain = endFailDomain.Name
	}
	if matchedPath.Domain == process.DomainNameRepeat {
		matchedPath.Domain = matchContext.Session.CurrentDomain
		matchedPath.Branch = matchContext.Session.CurrentBranch
	}
	slog.Info(fmt.Sprintf("sessionId [%v]: expected jump to [%v] actual jump to domain [%v] branch [%v]",
		matchContext.Session.SessionId, nextDomain, matchedPath.Domain, matchedPath.Branch),
		slog.Any("traceId", ctx.Value(traceid.TraceId{})))
	return makeAnswer(ctx, matchContext)
}

func makeRespond(matchContext *matcher.Context, answerDTO AnswerDTO, intentionRules []process.IntentionRule) *InteractiveRespond {
	interactiveRespond := &InteractiveRespond{}
	interactiveRespond.SessionId = matchContext.Session.SessionId
	interactiveRespond.Hits.Sentence = matchContext.Sentence
	if len(matchContext.Segments) == 0 {
		interactiveRespond.Hits.Segments = make([]string, 0)
	} else {
		interactiveRespond.Hits.Segments = matchContext.Segments
	}
	interactiveRespond.Hits.HitPaths = convertMatchedPathListToHitPathDTOList(matchContext.MatchedPaths)
	interactiveRespond.Answer.Text = replaceVariables(answerDTO.Text, matchContext.Session.Variables)
	interactiveRespond.Answer.Audio = answerDTO.Audio
	interactiveRespond.Intentions = convertIntentionRuleListToIntentionDTOList(intentionRules)
	return interactiveRespond
}

func convertMatchedPathListToHitPathDTOList(matchedPathList []matcher.MatchedPath) []HitPathDTO {
	if len(matchedPathList) == 0 {
		return make([]HitPathDTO, 0)
	}
	hitPathDTOList := make([]HitPathDTO, len(matchedPathList))
	for i, matchedPath := range matchedPathList {
		hitPathDTOList[i] = convertMatchedPathToHitPathDTO(matchedPath)
	}
	return hitPathDTOList
}

func convertMatchedPathToHitPathDTO(matchedPath matcher.MatchedPath) HitPathDTO {
	hitPathDTO := HitPathDTO{
		Domain: matchedPath.Domain,
		Branch: matchedPath.Branch,
	}
	if len(matchedPath.MatchedWords) == 0 {
		hitPathDTO.MatchedWords = make([]string, 0)
	} else {
		hitPathDTO.MatchedWords = matchedPath.MatchedWords
	}
	return hitPathDTO
}

func replaceVariables(text string, variables map[string]string) string {
	if len(variables) == 0 {
		return text
	}
	replacements := make([]string, 0)
	for code, value := range variables {
		replacements = append(replacements, code, value)
	}
	return strings.NewReplacer(replacements...).Replace(text)
}

func convertIntentionRuleListToIntentionDTOList(intentionRules []process.IntentionRule) []IntentionDTO {
	if len(intentionRules) == 0 {
		return make([]IntentionDTO, 0)
	}
	intentionDTOList := make([]IntentionDTO, len(intentionRules))
	for i, intentionRule := range intentionRules {
		intentionDTOList[i] = convertIntentionRuleToIntentionDTO(intentionRule)
	}
	return intentionDTOList
}

func convertIntentionRuleToIntentionDTO(intentionRule process.IntentionRule) IntentionDTO {
	return IntentionDTO{
		Code:        intentionRule.Code,
		DisplayName: intentionRule.DisplayName,
		Reason:      intentionRule.Reason,
	}
}
