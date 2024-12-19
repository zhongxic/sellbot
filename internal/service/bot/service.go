package bot

type Service interface {
	Prologue(prologueDTO *PrologueDTO) (*InteractiveRespond, error)
}

type ServiceImpl struct {
}

func NewService() Service {
	return &ServiceImpl{}
}
