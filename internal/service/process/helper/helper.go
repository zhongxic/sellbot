package helper

import (
	"fmt"
	"github.com/zhongxic/sellbot/internal/service/process"
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
		if domain.Type == process.DomainTypeStart {
			domains = append(domains, domain)
		}
	}
	if len(domains) != 1 {
		return process.Domain{}, fmt.Errorf("expected one start domain but found [%d]", len(domains))
	}
	return domains[0], nil
}
