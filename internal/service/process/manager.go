package process

type Manager struct {
	TestLoader    Loader
	ReleaseLoader Loader
}

func (manager *Manager) Load(processId string, test bool) (*Process, error) {
	if test {
		return manager.TestLoader.Load(processId)
	} else {
		return manager.ReleaseLoader.Load(processId)
	}
}
