package process

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Loader interface {
	Load(processId string) (*Process, error)
	LastModified(processId string) (time.Time, error)
}

func NewFileLoader(dir string) Loader {
	return &fileLoader{dir}
}

type fileLoader struct {
	dir string
}

func (loader *fileLoader) Load(processId string) (*Process, error) {
	path := filepath.Join(loader.dir, processId, fmt.Sprintf("%v.json", processId))
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	process := &Process{}
	err = json.Unmarshal(file, process)
	if err != nil {
		return nil, err
	}
	return process, nil
}

func (loader *fileLoader) LastModified(processId string) (time.Time, error) {
	path := filepath.Join(loader.dir, processId, fmt.Sprintf("%v.json", processId))
	stat, err := os.Stat(path)
	if err != nil {
		return time.Now(), err
	}
	return stat.ModTime(), nil
}
