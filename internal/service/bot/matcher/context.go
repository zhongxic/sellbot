package matcher

import (
	"errors"

	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/service/session"
)

type MatchedPath struct {
	Domain         string
	Branch         string
	DomainType     string
	DomainCategory string
	MatchedWords   []string
}

type Context struct {
	Session      *session.Session
	Process      *process.Process
	Sentence     string
	Segments     []string
	Silence      bool
	Interruption int
	MatchedPaths []MatchedPath
}

func (c *Context) AddMatchedPath(matchedPath MatchedPath) {
	c.MatchedPaths = append(c.MatchedPaths, matchedPath)
	// TODO update session stat
}

func (c *Context) GetLastMatchedPath() (MatchedPath, error) {
	if len(c.MatchedPaths) == 0 {
		return MatchedPath{}, errors.New("no matched path in context")
	}
	return c.MatchedPaths[len(c.MatchedPaths)-1], nil
}

func NewContext(session *session.Session, process *process.Process) *Context {
	return &Context{
		Session: session,
		Process: process,
	}
}
