package process

import (
	"fmt"
	"slices"
	"sort"
)

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
	if len(h.process.Domains) != 0 {
		for _, domain := range h.process.Domains {
			if domain.Category == DomainCategoryMainProcess && domain.Type == DomainTypeStart {
				domains = append(domains, domain)
			}
		}
	}
	if len(domains) != 1 {
		return Domain{}, fmt.Errorf("process [%v]: expected one start domain but found [%v]", h.process.Id, len(domains))
	}
	return domains[0], nil
}

func (h *Helper) GetSilenceDomain() (Domain, error) {
	domains := make([]Domain, 0)
	if len(h.process.Domains) != 0 {
		for _, domain := range h.process.Domains {
			if domain.Category == DomainCategorySilence && domain.Type == DomainTypeNormal {
				domains = append(domains, domain)
			}
		}
	}
	if len(domains) != 1 {
		return Domain{}, fmt.Errorf("process [%v]: expected one silence domain but found [%v]", h.process.Id, len(domains))
	}
	return domains[0], nil
}

func (h *Helper) GetCommonDialog(dialogType DomainType) (Domain, error) {
	domains := make([]Domain, 0)
	if len(h.process.Domains) != 0 {
		for _, domain := range h.process.Domains {
			if domain.Category == DomainCategoryCommonDialog && domain.Type == dialogType {
				domains = append(domains, domain)
			}
		}
	}
	if len(domains) != 1 {
		return Domain{}, fmt.Errorf("process [%v]: expected one [%v] common dialog but found [%d]",
			h.process.Id, dialogType, len(domains))
	}
	return domains[0], nil
}

func (h *Helper) GetBusinessQADomain() (Domain, error) {
	domains := make([]Domain, 0)
	if len(h.process.Domains) != 0 {
		for _, domain := range h.process.Domains {
			if domain.Category == DomainCategoryBusinessQA && domain.Type == DomainTypeNormal {
				domains = append(domains, domain)
			}
		}
	}
	if len(domains) != 1 {
		return Domain{}, fmt.Errorf("process [%v]: expected one business QA domain but found [%v]", h.process.Id, domains)
	}
	return domains[0], nil
}

func (h *Helper) GetBranch(domainName, branchName string) (Branch, error) {
	if len(h.process.Domains) == 0 {
		return Branch{}, fmt.Errorf("process [%v]: get branch failed due to empty domains", h.process.Id)
	}
	domain, ok := h.process.Domains[domainName]
	if !ok {
		return Branch{}, fmt.Errorf("process [%v]: get branch failed due to domain [%s] not found", h.process.Id, domainName)
	}
	branch, ok := domain.Branches[branchName]
	if !ok {
		return Branch{}, fmt.Errorf("process [%v]: get branch failed due to branch [%s] not found in domain [%v]",
			h.process.Id, branchName, domainName)
	}
	return branch, nil
}

func (h *Helper) GetDomainPositiveBranch(domainName string) (Branch, error) {
	if len(h.process.Domains) == 0 {
		return Branch{}, fmt.Errorf("process [%v]: get positive branch failed due to empty domains", h.process.Id)
	}
	domain, ok := h.process.Domains[domainName]
	if !ok {
		return Branch{}, fmt.Errorf("process [%v]: get positive branch failed due to domain [%s] not found", h.process.Id, domainName)
	}
	branches := make([]Branch, 0)
	if len(domain.Branches) != 0 {
		for _, branch := range domain.Branches {
			if branch.Semantic == BranchSemanticPositive {
				branches = append(branches, branch)
			}
		}
	}
	if len(branches) != 1 {
		return Branch{}, fmt.Errorf("process [%v]: expected one positive branch in domain [%v] but found [%v]",
			h.process.Id, domainName, len(branches))
	}
	return branches[0], nil
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

func (h *Helper) GetGlobalKeywords() ([]string, error) {
	keywords := make([]string, 0)
	for _, dialog := range DomainTypeDialogMatchOrders {
		domain, err := h.GetCommonDialog(dialog)
		if err != nil {
			return nil, fmt.Errorf("process [%v]: get global keywords failed due to get common dialog: %v", h.process.Id, err)
		}
		keywords = append(keywords, h.GetDomainKeywords(domain.Name)...)
	}
	domain, err := h.GetBusinessQADomain()
	if err != nil {
		return nil, fmt.Errorf("process [%v]: get global keywords failed due to get business QA: %v", h.process.Id, err)
	}
	keywords = append(keywords, h.GetDomainKeywords(domain.Name)...)
	return keywords, nil
}

func (h *Helper) GetMergeOrderedMatchPaths(domainName string) ([]MatchPath, error) {
	matchPaths := make([]MatchPath, 0)
	domain, err := h.GetDomain(domainName)
	if err != nil {
		return nil, fmt.Errorf("process [%v]: get domain [%v] match paths failed due to domain: %v",
			h.process.Id, domainName, err)
	}
	if len(domain.MatchOrders) > 0 {
		// user defined domain matcher order has the highest priority.
		matchPaths = append(matchPaths, domain.MatchOrders...)
	}
	// then common dialogs.
	for _, dialogType := range DomainTypeDialogMatchOrders {
		dialog, err := h.GetCommonDialog(dialogType)
		if err != nil {
			return nil, fmt.Errorf("process [%v]: get domain [%v] match paths failed due to dialog: %v",
				h.process.Id, domainName, err)
		}
		matchPaths = appendMatchPaths(matchPaths, dialog)
	}
	// next is others branches in this domain.
	matchPaths = appendMatchPaths(matchPaths, domain)
	// branches in business qa have the lowest priority.
	qaDomain, err := h.GetBusinessQADomain()
	if err != nil {
		return nil, fmt.Errorf("process [%v]: get domain [%v] match paths failed due to business QA: %v",
			h.process.Id, domainName, err)
	}
	matchPaths = appendMatchPaths(matchPaths, qaDomain)
	return matchPaths, nil
}

type branchSlice []Branch

func (s branchSlice) Len() int {
	return len(s)
}

func (s branchSlice) Less(i, j int) bool {
	return s[i].Order < s[j].Order
}

func (s branchSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func appendMatchPaths(matchPaths []MatchPath, domain Domain) []MatchPath {
	branches := make(branchSlice, 0)
	if len(domain.Branches) != 0 {
		for _, branch := range domain.Branches {
			branches = append(branches, branch)
		}
	}
	// consider business qa order.
	sort.Sort(branches)
	for _, branch := range branches {
		matchPath := MatchPath{DomainName: domain.Name, BranchName: branch.Name}
		if !slices.Contains(matchPaths, matchPath) {
			matchPaths = append(matchPaths, matchPath)
		}
	}
	return matchPaths
}

func NewHelper(process *Process) *Helper {
	return &Helper{process: process}
}
