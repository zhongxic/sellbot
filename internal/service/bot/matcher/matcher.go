package matcher

type Matcher interface {
	// Match find matched branches in process use this context
	// and abort match if return false
	Match(matchContext *Context) bool
}

func Match(matchContext *Context) {
	// TODO implement match
}
