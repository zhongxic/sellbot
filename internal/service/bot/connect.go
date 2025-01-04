package bot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/zhongxic/sellbot/internal/traceid"
)

func (s *serviceImpl) Connect(ctx context.Context, connectDTO *ConnectDTO) (*ConnectRespond, error) {
	slog.Info("start process connect", "traceId", ctx.Value(traceid.TraceId{}))
	currentSession, err := s.retrieveSession(connectDTO.SessionId)
	if err != nil {
		return nil, fmt.Errorf("retrieve session failed: %w", err)
	}
	tokenizer, err := s.retrieveTokenizer(connectDTO.SessionId)
	if err != nil {
		return nil, fmt.Errorf("retrieve tokenizer failed: %w", err)
	}
	currentSession.CallAnswerTime = time.Now()
	s.storeSession(currentSession.Id, currentSession)
	s.storeTokenizer(currentSession.Id, tokenizer)
	return &ConnectRespond{SessionId: currentSession.Id, AnswerTime: currentSession.CallAnswerTime.Format(time.DateTime)}, nil
}
