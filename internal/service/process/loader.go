package process

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zhongxic/sellbot/pkg/cache"
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
	lastModified, err := loader.LastModified(processId)
	if err != nil {
		return nil, err
	}
	process.lastModified = lastModified
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

type cachedLoader struct {
	rawLoader Loader
	storage   cache.Cache[string, *Process]
}

func (loader *cachedLoader) Load(processId string) (*Process, error) {
	process, exist := loader.storage.Get(processId)
	if exist {
		lastModified, err := loader.rawLoader.LastModified(processId)
		if err != nil || !process.lastModified.Equal(lastModified) {
			exist = false
		}
	}
	if !exist {
		process, err := loader.rawLoader.Load(processId)
		if err != nil {
			return nil, err
		}
		loader.storage.Set(processId, process)
		return process, nil
	}
	return process, nil
}

func (loader *cachedLoader) LastModified(processId string) (time.Time, error) {
	return loader.rawLoader.LastModified(processId)
}

func NewCachedLoader(rawLoader Loader, storage cache.Cache[string, *Process]) Loader {
	return &cachedLoader{
		rawLoader: rawLoader,
		storage:   storage,
	}
}
