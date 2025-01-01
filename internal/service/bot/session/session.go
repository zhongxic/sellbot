package session

import "github.com/google/uuid"

type StatPath struct {
}

type Session struct {
	SessionId             string
	ProcessId             string
	Variables             map[string]string
	Test                  bool
	CurrentDomain         string
	CurrentBranch         string
	LastMainProcessDomain string
	LastMainProcessBranch string
	ConversationCount     int
	DomainBranchHitCount  map[string]map[string]int
}

func New() *Session {
	return &Session{
		SessionId:            uuid.New().String(),
		DomainBranchHitCount: make(map[string]map[string]int),
	}
}

func (s *Session) GetDomainBranchHitCount(domainName, branchName string) int {
	if len(s.DomainBranchHitCount) == 0 {
		return 0
	}
	branchHitCount, ok := s.DomainBranchHitCount[domainName]
	if !ok || len(branchHitCount) == 0 {
		return 0
	}
	return branchHitCount[branchName]
}

func (s *Session) UpdateStat(statPath []StatPath) {
	// TODO impl-me update session stat
}
