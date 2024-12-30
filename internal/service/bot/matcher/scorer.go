package matcher

import (
	"context"

	"github.com/zhongxic/sellbot/internal/service/process"
)

type similarity struct {
	score   int
	matches []string
}

func (s similarity) isMatched() bool {
	return s.score > 0
}

func (s similarity) isBetterThan(other similarity) bool {
	return s.score > other.score
}

func score(ctx context.Context, text string, segments []string, Keywords process.Keywords) similarity {
	// TODO: implement the scoring logic
	return similarity{}
}
