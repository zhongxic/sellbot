package process

import (
	"fmt"
)

const (
	DefaultIntentionName   = "default"
	DefaultIntentionReason = "default intention"
)

type Helper struct {
	hold *Process
}

func NewHelper(hold *Process) *Helper {
	return &Helper{hold: hold}
}

func (h *Helper) GetDefaultIntentionRule() IntentionRule {
	return IntentionRule{
		Code:        h.hold.Intentions.DefaultIntention,
		DisplayName: DefaultIntentionName,
		Reason:      DefaultIntentionReason,
	}
}

func (h *Helper) GetDomain(domainName string) (Domain, error) {
	if len(h.hold.Domains) != 0 {
		if domain, ok := h.hold.Domains[domainName]; ok {
			return domain, nil
		}
	}
	return Domain{}, fmt.Errorf("process [%v]: domain [%v] not found", h.hold.Id, domainName)
}

func (h *Helper) GetStartDomain() (Domain, error) {
	domains := make([]Domain, 0)
	for _, domain := range h.hold.Domains {
		if domain.Category == DomainCategoryMainProcess && domain.Type == DomainTypeStart {
			domains = append(domains, domain)
		}
	}
	if len(domains) != 1 {
		return Domain{}, fmt.Errorf("process [%v]: expected one start domain but found [%v]", h.hold.Id, len(domains))
	}
	return domains[0], nil
}

func (h *Helper) GetSilenceDomain() (Domain, error) {
	// TODO impl-me get silence domain
	return Domain{}, nil
}

func (h *Helper) GetCommonDialogDomain(domainDialogType string) (Domain, error) {
	if len(h.hold.Domains) == 0 {
		return Domain{}, fmt.Errorf("process [%v]: empty domains", h.hold.Id)
	}
	domains := make([]Domain, 0)
	for _, domain := range h.hold.Domains {
		if domain.Category == DomainCategoryCommonDialog && domain.Type == domainDialogType {
			domains = append(domains, domain)
		}
	}
	if len(domains) != 1 {
		return Domain{}, fmt.Errorf("process [%v]: expected one [%v] common dialog but found [%d]",
			h.hold.Id, domainDialogType, len(domains))
	}
	return domains[0], nil
}

func (h *Helper) GetBranch(domainName, branchName string) (Branch, error) {
	if len(h.hold.Domains) == 0 {
		return Branch{}, fmt.Errorf("process [%v]: empty domains", h.hold.Id)
	}
	domain, ok := h.hold.Domains[domainName]
	if !ok {
		return Branch{}, fmt.Errorf("process [%v]: domain [%s] not found", h.hold.Id, domainName)
	}
	branch, ok := domain.Branches[branchName]
	if !ok {
		return Branch{}, fmt.Errorf("process [%v]: branch [%s] not found", h.hold.Id, branchName)
	}
	return branch, nil
}

func (h *Helper) GetDomainSemanticBranch(domainName, semantic string) (Branch, error) {
	// TODO impl-me get domain semantic branch
	return Branch{}, nil
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

func (h *Helper) GetMergeOrderedMatchPaths(domainName string) ([]MatchPath, error) {
	// TODO impl-me get merged match order
	return make([]MatchPath, 0), nil
}
