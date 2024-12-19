package bot

import "github.com/zhongxic/sellbot/internal/service/process"

type Service interface {
	Prologue(prologueDTO *PrologueDTO) (*InteractiveRespond, error)
}

type serviceImpl struct {
	processManager *process.Manager
}

type Options struct {
	TestProcessDir    string
	ReleaseProcessDir string
}

func NewService(options Options) Service {
	testLoader := process.NewFileLoader(options.TestProcessDir)
	releaseLoader := process.NewFileLoader(options.ReleaseProcessDir)
	return &serviceImpl{
		processManager: process.NewManager(testLoader, releaseLoader),
	}
}
