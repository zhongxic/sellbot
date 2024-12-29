package matcher

import (
	"context"
)

type Matcher interface {
	// Match find matched branches in process use this matchContext
	// and abort match if return false
	Match(ctx context.Context, matchContext *Context) (abort bool, err error)
}
