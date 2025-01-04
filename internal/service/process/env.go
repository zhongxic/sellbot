package process

import (
	"slices"
	"time"
)

type IntentionAnalyzeEnv struct {
	Sentence string
	Segments []string
	Status   Status `expr:"status"`
}

type Status struct {
	PreviousMainProcessDomain string
	ConversationCount         int
	PositiveCount             int
	NegativeCount             int
	RefusedCount              int
	BusinessQACount           int
	SilenceCount              int
	MissMatchCount            int
	CallAnswerTime            time.Time
	PassedDomains             []string
	DomainBranchHitCount      map[string]map[string]int
}

func (s Status) PassedDomain(domainName string) bool {
	if len(s.PassedDomains) == 0 {
		return false
	}
	return slices.Contains(s.PassedDomains, domainName)
}

func (s Status) GetDomainBranchHitCount(domainName, branchName string) int {
	branchHitCount := s.DomainBranchHitCount[domainName]
	if len(branchHitCount) == 0 {
		return 0
	}
	return branchHitCount[branchName]
}

func (s Status) AnswerSecondsCompareTo(sec int) int {
	if s.CallAnswerTime.IsZero() {
		return -1
	}
	return int(time.Since(s.CallAnswerTime).Seconds()) - sec
}
