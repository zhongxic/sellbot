package matcher

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/service/process/helper"
	"github.com/zhongxic/sellbot/internal/traceid"
)

type Matcher interface {
	// Match find matched branches in process use this matchContext
	// and abort match if return false
	Match(ctx context.Context, matchContext *Context) (abort bool, err error)
}

type OutOfMaxRoundsMatcher struct {
}

func (matcher *OutOfMaxRoundsMatcher) Match(ctx context.Context, matchContext *Context) (bool, error) {
	if matchContext.Session.ConversationCount >= matchContext.Process.Options.MaxRounds {
		slog.Info(fmt.Sprintf("sessionId [%v]: OutOfMaxRoundsMatcher detect conversation count out of max rounds [%v]",
			matchContext.Session.SessionId, matchContext.Process.Options.MaxRounds),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		processHelper := helper.New(matchContext.Process)
		domain, err := processHelper.FindCommonDialogDomain(process.DomainTypeDialogEndExceed)
		if err != nil {
			return true, fmt.Errorf("OutOfMaxRoundsMatcher find common dialog domain failed: %w", err)
		}
		matchedPath := MatchedPath{Domain: domain.Name, Branch: process.BranchNameEnter}
		slog.Info(fmt.Sprintf("sessionId [%v]: OutOfMaxRoundsMatcher matched domain [%v] branch [%v]",
			matchContext.Session.SessionId, matchedPath.Domain, matchedPath.Branch),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		matchContext.AddMatchedPath(matchedPath)
		return true, nil
	}
	return false, nil
}

type ForceInterruptionMatcher struct {
}

func (matcher *ForceInterruptionMatcher) Match(ctx context.Context, matchContext *Context) (bool, error) {
	if matchContext.Interruption == process.InterruptionTypeForce {
		slog.Info(fmt.Sprintf("sessionId [%v]: ForceInterruptionMatcher detect force interruption", matchContext.Session.SessionId),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		processHelper := helper.New(matchContext.Process)
		domain, err := processHelper.GetForceInterruptionJumpToDomain()
		if err != nil {
			return true, fmt.Errorf("ForceInterruptionMatcher get force interrupt jump to domain failed: %w", err)
		}
		matchedPath := MatchedPath{Domain: domain.Name, Branch: process.BranchNameEnter}
		slog.Info(fmt.Sprintf("sessionId [%v]: ForceInterruptionMatcher matched domain [%v] branch [%v]",
			matchContext.Session.SessionId, matchedPath.Domain, matchedPath.Branch),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		matchContext.AddMatchedPath(matchedPath)
		return true, nil
	}
	return false, nil
}
