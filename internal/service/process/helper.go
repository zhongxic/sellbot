package process

import "fmt"

const (
	DefaultIntentionName   = "default"
	DefaultIntentionReason = "default intention"
)

type Helper struct {
	process *Process
}

func (h *Helper) GetDefaultIntentionRule() IntentionRule {
	return IntentionRule{
		Code:        h.process.Intentions.DefaultIntention,
		DisplayName: DefaultIntentionName,
		Reason:      DefaultIntentionReason,
	}
}

func (h *Helper) GetDomain(domainName string) (Domain, error) {
	if len(h.process.Domains) != 0 {
		if domain, ok := h.process.Domains[domainName]; ok {
			return domain, nil
		}
	}
	return Domain{}, fmt.Errorf("process [%v]: domain [%v] not found", h.process.Id, domainName)
}

func (h *Helper) GetStartDomain() (Domain, error) {
	domains := make([]Domain, 0)
	for _, domain := range h.process.Domains {
		if domain.Category == DomainCategoryMainProcess && domain.Type == DomainTypeStart {
			domains = append(domains, domain)
		}
	}
	if len(domains) != 1 {
		return Domain{}, fmt.Errorf("process [%v]: expected one start domain but found [%v]", h.process.Id, len(domains))
	}
	return domains[0], nil
}

func (h *Helper) GetSilenceDomain() (Domain, error) {
	// TODO impl-me get silence domain
	return Domain{}, nil
}

func (h *Helper) GetCommonDialog(dialogType DomainType) (Domain, error) {
	if len(h.process.Domains) == 0 {
		return Domain{}, fmt.Errorf("process [%v]: empty domains", h.process.Id)
	}
	domains := make([]Domain, 0)
	for _, domain := range h.process.Domains {
		if domain.Category == DomainCategoryCommonDialog && domain.Type == dialogType {
			domains = append(domains, domain)
		}
	}
	if len(domains) != 1 {
		return Domain{}, fmt.Errorf("process [%v]: expected one [%v] common dialog but found [%d]",
			h.process.Id, dialogType, len(domains))
	}
	return domains[0], nil
}

func (h *Helper) GetBranch(domainName, branchName string) (Branch, error) {
	if len(h.process.Domains) == 0 {
		return Branch{}, fmt.Errorf("process [%v]: empty domains", h.process.Id)
	}
	domain, ok := h.process.Domains[domainName]
	if !ok {
		return Branch{}, fmt.Errorf("process [%v]: domain [%s] not found", h.process.Id, domainName)
	}
	branch, ok := domain.Branches[branchName]
	if !ok {
		return Branch{}, fmt.Errorf("process [%v]: branch [%s] not found", h.process.Id, branchName)
	}
	return branch, nil
}

func (h *Helper) GetDomainSemanticBranch(domainName string, semantic BranchSemantic) (Branch, error) {
	// TODO impl-me get domain semantic branch
	return Branch{}, nil
}

func (h *Helper) GetDomainKeywords(domainName string) []string {
	keywords := make([]string, 0)
	if len(h.process.Domains) == 0 {
		return keywords
	}
	domain, ok := h.process.Domains[domainName]
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
	if len(h.process.Domains) == 0 {
		return keywords
	}
	domain, ok := h.process.Domains[domainName]
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

func NewHelper(process *Process) *Helper {
	return &Helper{process: process}
}
