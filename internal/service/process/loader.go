package process

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/expr-lang/expr"
	"github.com/zhongxic/sellbot/pkg/cache"
)

type Loader interface {
	Load(processId string) (*Process, error)
	LastModified(processId string) (time.Time, error)
}

type fileLoader struct {
	dir string
}

func (loader *fileLoader) Load(processId string) (*Process, error) {
	path := filepath.Join(loader.dir, processId, fmt.Sprintf("%v.json", processId))
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("process [%v] read file failed: %w", processId, err)
	}
	process := &Process{}
	err = json.Unmarshal(file, process)
	if err != nil {
		return nil, fmt.Errorf("process [%v] unmarshal failed: %w", processId, err)
	}
	lastModified, err := loader.LastModified(processId)
	if err != nil {
		return nil, fmt.Errorf("process [%v] read stat failed: %w", processId, err)
	}
	// compile intention rules expression.
	if err := compileIntentionRulesExpression(process); err != nil {
		return nil, fmt.Errorf("process [%v] compile intention rules failed: %w", processId, err)
	}
	process.lastModified = lastModified
	return process, nil
}

func compileIntentionRulesExpression(process *Process) error {
	size := len(process.Intentions.IntentionRules)
	if size == 0 {
		return nil
	}
	env := IntentionAnalyzeEnv{}
	for i := 0; i < size; i++ {
		rule := &process.Intentions.IntentionRules[i]
		if rule.Expression != "" {
			program, err := expr.Compile(rule.Expression, expr.Env(env), expr.AsBool())
			if err != nil {
				return fmt.Errorf("rule [%v] expression [%v]  compile failed: %w", rule.Code, rule.Expression, err)
			}
			rule.program = program
		}
	}
	return nil
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

func NewFileLoader(dir string) Loader {
	return &fileLoader{dir}
}

func NewCachedLoader(rawLoader Loader, storage cache.Cache[string, *Process]) Loader {
	return &cachedLoader{
		rawLoader: rawLoader,
		storage:   storage,
	}
}
