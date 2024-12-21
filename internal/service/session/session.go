package session

import "github.com/google/uuid"

type Session struct {
	SessionId             string
	ProcessId             string
	Variables             map[string]string
	Test                  bool
	CurrentDomain         string
	CurrentBranch         string
	LastMainProcessDomain string
	LastMainProcessBranch string
}

func New() *Session {
	return &Session{SessionId: uuid.New().String()}
}
