package helper

import (
	"fmt"

	"github.com/zhongxic/sellbot/internal/service/process"
)

const (
	DefaultIntentionName   = "default"
	DefaultIntentionReason = "default intention"
)

type Helper struct {
	hold *process.Process
}

func New(hold *process.Process) *Helper {
	return &Helper{hold: hold}
}

func (h *Helper) GetDefaultIntentionRule() process.IntentionRule {
	return process.IntentionRule{
		Code:        h.hold.Intentions.DefaultIntention,
		DisplayName: DefaultIntentionName,
		Reason:      DefaultIntentionReason,
	}
}

func (h *Helper) GetDomain(domainName string) (process.Domain, error) {
	if len(h.hold.Domains) != 0 {
		if domain, ok := h.hold.Domains[domainName]; ok {
			return domain, nil
		}
	}
	return process.Domain{}, fmt.Errorf("process [%v]: domain [%v] not found", h.hold.Id, domainName)
}

func (h *Helper) GetStartDomain() (process.Domain, error) {
	domains := make([]process.Domain, 0)
	for _, domain := range h.hold.Domains {
		if domain.Category == process.DomainCategoryMainProcess && domain.Type == process.DomainTypeStart {
			domains = append(domains, domain)
		}
	}
	if len(domains) != 1 {
		return process.Domain{}, fmt.Errorf("process [%v]: expected one start domain but found [%v]", h.hold.Id, len(domains))
	}
	return domains[0], nil
}

func (h *Helper) GetSilenceDomain() (process.Domain, error) {
	// TODO impl-me get silence domain
	return process.Domain{}, nil
}

func (h *Helper) GetCommonDialogDomain(domainDialogType string) (process.Domain, error) {
	if len(h.hold.Domains) == 0 {
		return process.Domain{}, fmt.Errorf("process [%v]: empty domains", h.hold.Id)
	}
	domains := make([]process.Domain, 0)
	for _, domain := range h.hold.Domains {
		if domain.Category == process.DomainCategoryCommonDialog && domain.Type == domainDialogType {
			domains = append(domains, domain)
		}
	}
	if len(domains) != 1 {
		return process.Domain{}, fmt.Errorf("process [%v]: expected one [%v] common dialog but found [%d]",
			h.hold.Id, domainDialogType, len(domains))
	}
	return domains[0], nil
}

func (h *Helper) GetBranch(domainName, branchName string) (process.Branch, error) {
	if len(h.hold.Domains) == 0 {
		return process.Branch{}, fmt.Errorf("process [%v]: empty domains", h.hold.Id)
	}
	domain, ok := h.hold.Domains[domainName]
	if !ok {
		return process.Branch{}, fmt.Errorf("process [%v]: domain [%s] not found", h.hold.Id, domainName)
	}
	branch, ok := domain.Branches[branchName]
	if !ok {
		return process.Branch{}, fmt.Errorf("process [%v]: branch [%s] not found", h.hold.Id, branchName)
	}
	return branch, nil
}

func (h *Helper) GetDomainSemanticBranch(domainName, semantic string) (process.Branch, error) {
	// TODO impl-me get domain semantic branch
	return process.Branch{}, nil
}

func (h *Helper) GetDomainKeywords(domainName string) []string {
	keywords := make([]string, 0)
	if len(h.hold.Domains) == 0 {
		return keywords
	}
	domain, ok := h.hold.Domains[domainName]
	if !ok || len(domain.Branches) == 0 {
		return keywords
	}
	for branchName := range domain.Branches {
		flatMapKeywords := h.GetBranchKeywords(domainName, branchName)
		keywords = append(keywords, flatMapKeywords...)
	}
	return keywords
}

func (h *Helper) GetBranchKeywords(domainName, branchName string) []string {
	keywords := make([]string, 0)
	if len(h.hold.Domains) == 0 {
		return keywords
	}
	domain, ok := h.hold.Domains[domainName]
	if !ok || len(domain.Branches) == 0 {
		return keywords
	}
	branch, ok := domain.Branches[branchName]
	if !ok {
		return keywords
	}
	if len(branch.Keywords.Simple) > 0 {
		keywords = append(keywords, branch.Keywords.Simple...)
	}
	if len(branch.Keywords.Combination) > 0 {
		for _, words := range branch.Keywords.Combination {
			if len(words) > 0 {
				keywords = append(keywords, words...)
			}
		}
	}
	if len(branch.Keywords.Exact) > 0 {
		keywords = append(keywords, branch.Keywords.Exact...)
	}
	return keywords
}

func (h *Helper) GetGlobalKeywords() []string {
	// TODO impl-me load global keywords
	return make([]string, 0)
}

func (h *Helper) GetMergeOrderedMatchPaths(domainName string) ([]process.MatchPath, error) {
	// TODO impl-me get merged match order
	return make([]process.MatchPath, 0), nil
}
