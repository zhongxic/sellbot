package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/zhongxic/sellbot/internal/service/bot/session"
	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/traceid"
)

func assembleIntentionAnalyzeEnv(sentence string, segments []string, sess *session.Session) process.IntentionAnalyzeEnv {
	return process.IntentionAnalyzeEnv{
		Sentence: sentence,
		Segments: segments,
		Status: process.Status{
			PreviousMainProcessDomain: sess.PreviousMainProcessDomain,
			ConversationCount:         sess.ConversationCount,
			PositiveCount:             sess.PositiveCount,
			NegativeCount:             sess.NegativeCount,
			RefusedCount:              sess.RefusedCount,
			BusinessQACount:           sess.BusinessQACount,
			SilenceCount:              sess.SilenceCount,
			MissMatchCount:            sess.MissMatchCount,
			CallAnswerTime:            sess.CallAnswerTime,
			PassedDomains:             sess.PassedDomains,
			DomainBranchHitCount:      sess.DomainBranchHitCount,
		},
	}
}

func analyzeIntention(ctx context.Context, env process.IntentionAnalyzeEnv, rules []process.IntentionRule) []process.IntentionRule {
	matchedRules := make([]process.IntentionRule, 0)
	if len(rules) == 0 {
		return matchedRules
	}
	for _, rule := range rules {
		hit, err := rule.IsHit(ctx, env)
		if err != nil {
			slog.Error(fmt.Sprintf("judge intention rule [%v] failed: %v", rule.Code, err),
				slog.Any("traceId", ctx.Value(traceid.TraceId{})))
			continue
		}
		if hit {
			matchedRules = append(matchedRules, rule)
		}
	}
	return matchedRules
}
