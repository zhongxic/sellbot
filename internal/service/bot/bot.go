package bot

import (
	"context"

	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/service/session"
)

type Service interface {
	Prologue(ctx context.Context, prologueDTO *PrologueDTO) (*InteractiveRespond, error)
}

type serviceImpl struct {
	options        Options
	processManager *process.Manager
}

func (s *serviceImpl) initSession(prologueDTO *PrologueDTO) *session.Session {
	sess := session.New()
	sess.ProcessId = prologueDTO.ProcessId
	sess.Variables = prologueDTO.Variables
	sess.Test = prologueDTO.Test
	// TODO add session to cache
	return sess
}

type Options struct {
	DictFile          string
	TestProcessDir    string
	ReleaseProcessDir string
}

func NewService(options Options, processManager *process.Manager) Service {
	return &serviceImpl{
		options:        options,
		processManager: processManager,
	}
}
