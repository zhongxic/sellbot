package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/zhongxic/sellbot/internal/traceid"
)

func (s *serviceImpl) Hangup(ctx context.Context, sessionIdDTO *SessionIdDTO) (*InteractiveRespond, error) {
	slog.Info("start process hold", "traceId", ctx.Value(traceid.TraceId{}))
	currentSession, err := s.retrieveSession(sessionIdDTO.SessionId)
	if err != nil {
		return nil, fmt.Errorf("retrieve session failed: %w", err)
	}
	loadedProcess, err := s.processManager.Load(currentSession.ProcessId, currentSession.Test)
	if err != nil {
		return nil, fmt.Errorf("load process failed: %w", err)
	}
	env := assembleIntentionAnalyzeEnv("", nil, currentSession)
	intentionRules := analyzeIntention(ctx, env, loadedProcess.Intentions.IntentionRules)
	interactiveRespond := &InteractiveRespond{
		SessionId:  currentSession.Id,
		Intentions: convertIntentionRuleListToIntentionDTOList(intentionRules),
	}
	s.invalidateSession(currentSession.Id)
	s.invalidateTokenizer(currentSession.Id)
	currentSession = nil
	return interactiveRespond, nil
}
