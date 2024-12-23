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

func (h *Helper) FindStartDomain() (process.Domain, error) {
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

func (h *Helper) GetBranch(domainName, branchName string) (process.Branch, error) {
	if len(h.hold.Domains) == 0 {
		return process.Branch{}, fmt.Errorf("process [%v]: empty domains", h.hold.Id)
	}
	domain, ok := h.hold.Domains[domainName]
	if !ok {
		return process.Branch{}, fmt.Errorf("process [%v]: domain [%s] not foun", h.hold.Id, domainName)
	}
	branch, ok := domain.Branches[branchName]
	if !ok {
		return process.Branch{}, fmt.Errorf("process [%v]: branch [%s] not found", h.hold.Id, branchName)
	}
	return branch, nil
}

func (h *Helper) FindCommonDialogDomain(domainDialogType string) (process.Domain, error) {
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
