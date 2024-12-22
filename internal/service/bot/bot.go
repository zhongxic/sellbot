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
	TestProcessDir    string
	ReleaseProcessDir string
}

func NewService(options Options) Service {
	// TODO add cached loader impl
	testLoader := process.NewFileLoader(options.TestProcessDir)
	releaseLoader := process.NewFileLoader(options.ReleaseProcessDir)
	return &serviceImpl{
		processManager: process.NewManager(testLoader, releaseLoader),
	}
}
