package session

import (
	"github.com/google/uuid"
	"github.com/zhongxic/sellbot/internal/service/process"
	"slices"
	"time"
)

type HitPathView struct {
	Domain         string
	Branch         string
	DomainType     process.DomainType
	DomainCategory process.DomainCategory
	BranchSemantic process.BranchSemantic
}

// Session is current session
//
// Notice that CurrentMainProcessDomain is located main process domain in current interaction before session update,
// PreviousMainProcessDomain is located main process domain in current interaction after session updated.
// Process matching should use CurrentMainProcessDomain, and intention analyze should use PreviousMainProcessDomain because
// session has been updated.
type Session struct {
	Id                        string
	ProcessId                 string
	Variables                 map[string]string
	Test                      bool
	CurrentDomain             string
	CurrentBranch             string
	CurrentMainProcessDomain  string
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

func (s *Session) UpdateStat(hitPaths []HitPathView) {
	if len(hitPaths) == 0 {
		return
	}
	if isHitPositiveBranch(hitPaths) {
		s.PositiveCount++
	}
	if isHitNegativeDomainOrBranch(hitPaths) {
		s.NegativeCount++
	}
	if isHitRefusedDomain(hitPaths) {
		s.RefusedCount++
	}
	if isHitBusinessQA(hitPaths) {
		s.BusinessQACount++
	}
	if isHitSilenceDomain(hitPaths) {
		s.SilenceCount++
	}
	if isHitMissMatchDomain(hitPaths) {
		s.MissMatchCount++
	}
	s.ConversationCount++
	lastMatchedPath := hitPaths[len(hitPaths)-1]
	s.CurrentDomain = lastMatchedPath.Domain
	s.CurrentBranch = lastMatchedPath.Branch
	if lastMatchedPath.DomainCategory == process.DomainCategoryMainProcess {
		s.PreviousMainProcessDomain = s.CurrentMainProcessDomain
		s.CurrentMainProcessDomain = lastMatchedPath.Domain
	}
	for _, path := range hitPaths {
		s.PassedDomains = append(s.PassedDomains, path.Domain)
		branchHitCount := computeIfAbsent[string, map[string]int](s.DomainBranchHitCount, path.Domain,
			func(key string) map[string]int {
				return make(map[string]int)
			})
		branchHitCount[path.Branch] = branchHitCount[path.Branch] + 1
	}
}

func computeIfAbsent[K comparable, V any](m map[K]V, key K, compute func(key K) V) (value V) {
	if m == nil {
		return
	}
	value, ok := m[key]
	if !ok {
		value = compute(key)
		m[key] = value
	}
	return
}

func isHitPositiveBranch(statPaths []HitPathView) bool {
	if len(statPaths) == 0 {
		return false
	}
	for _, path := range statPaths {
		if path.BranchSemantic == process.BranchSemanticPositive {
			return true
		}
	}
	return false
}

func isHitNegativeDomainOrBranch(statPaths []HitPathView) bool {
	if len(statPaths) == 0 {
		return false
	}
	for _, path := range statPaths {
		isNegative := slices.Contains(process.NegativeDomainTypes, path.DomainType) ||
			path.BranchSemantic == process.BranchSemanticNegative
		if isNegative {
			return true
		}
	}
	return false
}

func isHitRefusedDomain(statPaths []HitPathView) bool {
	if len(statPaths) == 0 {
		return false
	}
	for _, path := range statPaths {
		if path.DomainCategory == process.DomainCategoryCommonDialog &&
			path.DomainType == process.DomainTypeDialogRefused {
			return true
		}
	}
	return false
}

func isHitBusinessQA(statPaths []HitPathView) bool {
	if len(statPaths) == 0 {
		return false
	}
	for _, path := range statPaths {
		if path.DomainCategory == process.DomainCategoryBusinessQA {
			return true
		}
	}
	return false
}

func isHitSilenceDomain(statPaths []HitPathView) bool {
	if len(statPaths) == 0 {
		return false
	}
	for _, path := range statPaths {
		if path.DomainCategory == process.DomainCategorySilence {
			return true
		}
	}
	return false
}

func isHitMissMatchDomain(statPaths []HitPathView) bool {
	if len(statPaths) == 0 {
		return false
	}
	for _, path := range statPaths {
		if path.DomainCategory == process.DomainCategoryCommonDialog &&
			path.DomainType == process.DomainTypeDialogMissMatch {
			return true
		}
	}
	return false
}

func New() *Session {
	return &Session{
		Id:                   uuid.New().String(),
		DomainBranchHitCount: make(map[string]map[string]int),
		PassedDomains:        make([]string, 0),
	}
}
