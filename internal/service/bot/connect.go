package bot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/zhongxic/sellbot/internal/traceid"
)

func (s *serviceImpl) Connect(ctx context.Context, sessionIdDTO *SessionIdDTO) (*InteractiveRespond, error) {
	slog.Info("start process connect", "traceId", ctx.Value(traceid.TraceId{}))
	currentSession, err := s.retrieveSession(sessionIdDTO.SessionId)
	if err != nil {
		return nil, fmt.Errorf("retrieve session failed: %w", err)
	}
	tokenizer, err := s.retrieveTokenizer(sessionIdDTO.SessionId)
	if err != nil {
		return nil, fmt.Errorf("retrieve tokenizer failed: %w", err)
	}
	currentSession.CallAnswerTime = time.Now()
	s.storeSession(currentSession.Id, currentSession)
	s.storeTokenizer(currentSession.Id, tokenizer)
	return &InteractiveRespond{SessionId: currentSession.Id}, nil
}
