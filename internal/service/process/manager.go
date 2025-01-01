package process

type Manager struct {
	TestProcessLoader    Loader
	ReleaseProcessLoader Loader
}

func (m *Manager) Load(processId string, test bool) (*Process, error) {
	if test {
		return m.TestProcessLoader.Load(processId)
	}
	return m.ReleaseProcessLoader.Load(processId)
}
