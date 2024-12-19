package bot

type Service interface {
	Prologue(prologueDTO *PrologueDTO) *InteractiveRespond
}

type ServiceImpl struct {
}

func NewService() Service {
	return &ServiceImpl{}
}
