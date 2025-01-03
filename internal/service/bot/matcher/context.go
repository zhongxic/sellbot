package matcher

import (
	"errors"
	"fmt"

	"github.com/zhongxic/sellbot/internal/service/bot/session"
	"github.com/zhongxic/sellbot/internal/service/process"
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
	Silence      bool
	Interruption int
	MatchedPaths []MatchedPath
}

func (c *Context) AddMatchedPath(matchedPath MatchedPath) {
	c.MatchedPaths = append(c.MatchedPaths, matchedPath)
}

func (c *Context) GetLastMatchedPath() (MatchedPath, error) {
	if len(c.MatchedPaths) == 0 {
		return MatchedPath{}, errors.New("no matched path in context")
	}
	return c.MatchedPaths[len(c.MatchedPaths)-1], nil
}

func (c *Context) UpdateSessionStat() error {
	if len(c.MatchedPaths) == 0 {
		return nil
	}
	hitPaths := make([]session.HitPathView, 0)
	for _, matchedPath := range c.MatchedPaths {
		path, err := c.convertToHitPathView(matchedPath)
		if err != nil {
			return fmt.Errorf("convert to hit path failed: %w", err)
		}
		hitPaths = append(hitPaths, path)
	}
	c.Session.UpdateStat(hitPaths)
	return nil
}

func (c *Context) convertToHitPathView(matchedPath MatchedPath) (hitPath session.HitPathView, err error) {
	processHelper := process.NewHelper(c.Process)
	domain, err := processHelper.GetDomain(matchedPath.Domain)
	if err != nil {
		return session.HitPathView{}, err
	}
	branch, err := processHelper.GetBranch(matchedPath.Domain, matchedPath.Branch)
	if err != nil {
		return session.HitPathView{}, err
	}
	hitPath = session.HitPathView{
		Domain:         domain.Name,
		Branch:         branch.Name,
		DomainType:     domain.Type,
		DomainCategory: domain.Category,
		BranchSemantic: branch.Semantic,
	}
	return
}

func NewContext(session *session.Session, process *process.Process) *Context {
	return &Context{
		Session: session,
		Process: process,
	}
}
