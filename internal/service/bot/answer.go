package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/zhongxic/sellbot/internal/service/bot/matcher"
	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/traceid"
)

var maxAllowedJumpCount = 5

// autoJumpCountContextKey is a context key to retrieve auto jumped count in current process.
type autoJumpCountContextKey struct {
}

func makeAnswer(ctx context.Context, matchContext *matcher.Context) (AnswerDTO, error) {
	if matchContext == nil {
		return AnswerDTO{}, fmt.Errorf("make answer failed due to nil match context")
	}
	matchedPath, err := matchContext.GetLastMatchedPath()
	if err != nil {
		return AnswerDTO{}, fmt.Errorf("get last matched path failed: %w", err)
	}
	traceId := slog.Any("traceId", ctx.Value(traceid.TraceId{}))
	slog.Info(fmt.Sprintf("sessionId [%v]: current domain [%v] currnet mainProcessDomain [%v] matched domain [%v] branch [%v]",
		matchContext.Session.Id, matchContext.Session.CurrentDomain, matchContext.Session.CurrentMainProcessDomain,
		matchedPath.Domain, matchedPath.Branch), traceId)
	processHelper := process.NewHelper(matchContext.Process)
	domain, err := processHelper.GetDomain(matchedPath.Domain)
	if err != nil {
		return AnswerDTO{}, fmt.Errorf("get domain failed: %w", err)
	}
	branch, err := processHelper.GetBranch(matchedPath.Domain, matchedPath.Branch)
	if err != nil {
		return AnswerDTO{}, fmt.Errorf("get branch failed: %w", err)
	}
	if len(branch.Responses) == 0 {
		return autoJump(ctx, matchContext, branch.Next)
	}
	hitCount := matchContext.Session.GetDomainBranchHitCount(matchedPath.Domain, matchedPath.Branch)
	isExceed := hitCount >= len(branch.Responses) && domain.Category != process.DomainCategoryMainProcess
	slog.Info(fmt.Sprintf("sessionId [%v]: domain [%v] branch [%v] hitCount [%v] isExceed [%v]",
		matchContext.Session.Id, matchedPath.Domain, matchedPath.Branch, hitCount, isExceed), traceId)
	if isExceed {
		nextDomain := ""
		if branch.EnableExceedJump && branch.Next != "" {
			nextDomain = branch.Next
		}
		slog.Info(fmt.Sprintf("sessionId [%v]: jump to domain [%v] due to hitCount exceed",
			matchContext.Session.Id, nextDomain), traceId)
		return autoJump(ctx, matchContext, nextDomain)
	}
	ended, agent := false, false
	if domain.Type.IsEnded() {
		ended = true
	}
	if domain.Type == process.DomainTypeAgent {
		agent = true
	}
	response := branch.Responses[hitCount%len(branch.Responses)]
	if response.EnableAutoJump && response.Next != "" {
		slog.Info(fmt.Sprintf("sessionId [%v]: jump to domain [%v] due to domain [%v] branch [%v] auto jump enabled",
			matchContext.Session.Id, response.Next, matchedPath.Domain, matchedPath.Branch), traceId)
		return autoJump(ctx, matchContext, response.Next)
	}
	return AnswerDTO{Text: response.Text, Audio: response.Audio, Ended: ended, Agent: agent}, nil
}

func autoJump(ctx context.Context, matchContext *matcher.Context, nextDomain string) (AnswerDTO, error) {
	processHelper := process.NewHelper(matchContext.Process)
	// make sure not loop forever
	autoJumpCount := ctx.Value(autoJumpCountContextKey{})
	if autoJumpCount != nil && autoJumpCount.(int) >= maxAllowedJumpCount {
		return AnswerDTO{}, fmt.Errorf("max allowed jump count exceed")
	}
	var jumpContext context.Context
	if autoJumpCount == nil {
		jumpContext = context.WithValue(ctx, autoJumpCountContextKey{}, 0)
	} else {
		jumpContext = context.WithValue(ctx, autoJumpCountContextKey{}, autoJumpCount.(int)+1)
	}
	var matchedPath matcher.MatchedPath
	if nextDomain == "" {
		domain, err := processHelper.GetCommonDialog(process.DomainTypeDialogEndFail)
		if err != nil {
			return AnswerDTO{}, fmt.Errorf("get common dialog [%v] failed: %w", process.DomainTypeDialogEndFail, err)
		}
		matchedPath = matcher.MatchedPath{Domain: domain.Name, Branch: process.BranchNameEnter}
	} else if nextDomain == process.DomainNameRepeat {
		matchedPath = matcher.MatchedPath{Domain: matchContext.Session.CurrentDomain, Branch: matchContext.Session.CurrentBranch}
	} else {
		matchedPath = matcher.MatchedPath{Domain: nextDomain, Branch: process.BranchNameEnter}
	}
	slog.Info(fmt.Sprintf("sessionId [%v]: expected jump to [%v] actual jump to domain [%v] branch [%v]",
		matchContext.Session.Id, nextDomain, matchedPath.Domain, matchedPath.Branch),
		slog.Any("traceId", ctx.Value(traceid.TraceId{})))
	matchContext.AddMatchedPath(matchedPath)
	return makeAnswer(jumpContext, matchContext)
}
