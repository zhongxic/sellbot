package matcher

import (
	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/service/session"
)

type MatchedPath struct {
	Domain       string
	Branch       string
	MatchedWords []string
}

type Context struct {
	Session      *session.Session
	Process      *process.Process
	Sentence     string
	Segments     []string
	MatchedPaths []MatchedPath
}

func (c *Context) AddMatchedPath(matchedPath MatchedPath) {
	c.MatchedPaths = append(c.MatchedPaths, matchedPath)
	// TODO update session stat
}

func NewContext(session *session.Session, process *process.Process) *Context {
	return &Context{
		Session: session,
		Process: process,
	}
}
