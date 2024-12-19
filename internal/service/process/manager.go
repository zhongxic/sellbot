package process

type Manager struct {
	TestLoader    Loader
	ReleaseLoader Loader
}

func NewManager(testLoader Loader, releaseLoader Loader) *Manager {
	return &Manager{
		TestLoader:    testLoader,
		ReleaseLoader: releaseLoader,
	}
}

func (m *Manager) Load(processId string, test bool) (*Process, error) {
	if test {
		return m.TestLoader.Load(processId)
	} else {
		return m.ReleaseLoader.Load(processId)
	}
}
