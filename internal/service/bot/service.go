package bot

type Service interface {
	Prologue(prologueDTO *PrologueDTO) (*InteractiveRespond, error)
}

type serviceImpl struct {
}

func NewService() Service {
	return &serviceImpl{}
}
