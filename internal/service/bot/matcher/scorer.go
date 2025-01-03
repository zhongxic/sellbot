package matcher

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/traceid"
)

const (
	simpleKeywordBoost       = 1
	combinationKeywordsBoost = 2
	exactKeywordsBoost       = 100
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

func (s similarity) String() string {
	return fmt.Sprintf("{score: %v, matches: %v}", s.score, s.matches)
}

func score(ctx context.Context, text string, segments []string, Keywords process.Keywords) similarity {
	simpleKeywordsSim := scoreSimpleKeywords(segments, Keywords.Simple)
	combinationKeywordsSim := scoreCombinationKeywords(segments, Keywords.Combination)
	exactKeywordsSim := scoreExactKeywords(segments, Keywords.Exact)
	slog.Debug(fmt.Sprintf("text [%v] segments [%v]: "+
		"simple keywords similarity [%v], combination keywords similarity [%v], exact keywords similarity [%v]",
		text, segments, simpleKeywordsSim, combinationKeywordsSim, exactKeywordsSim),
		slog.Any("traceId", ctx.Value(traceid.TraceId{})))
	maxSim := similarity{}
	if simpleKeywordsSim.isBetterThan(maxSim) {
		maxSim = simpleKeywordsSim
	}
	if combinationKeywordsSim.isBetterThan(maxSim) {
		maxSim = combinationKeywordsSim
	}
	if exactKeywordsSim.isBetterThan(maxSim) {
		maxSim = exactKeywordsSim
	}
	slog.Debug(fmt.Sprintf("text [%v] segments [%v]: best matched keywords similarity is [%v]",
		text, segments, maxSim),
		slog.Any("traceId", ctx.Value(traceid.TraceId{})))
	return maxSim
}

func scoreSimpleKeywords(segments, simpleKeywords []string) similarity {
	if len(segments) == 0 || len(simpleKeywords) == 0 {
		return similarity{}
	}
	var matches []string
	for _, segment := range segments {
		if slices.Contains(simpleKeywords, segment) {
			matches = append(matches, segment)
		}
	}
	return similarity{
		score:   len(matches) * simpleKeywordBoost,
		matches: matches,
	}
}

func scoreCombinationKeywords(segments []string, combinationKeywords [][]string) similarity {
	maxSim := similarity{}
	if len(segments) == 0 || len(combinationKeywords) == 0 {
		return maxSim
	}
	for _, combinations := range combinationKeywords {
		if keywordsMatchAll(segments, combinations) {
			sim := similarity{
				score:   len(combinations) * combinationKeywordsBoost,
				matches: combinations,
			}
			if sim.isBetterThan(maxSim) {
				maxSim = sim
			}
		}
	}
	return maxSim
}

func keywordsMatchAll(segments, keywords []string) bool {
	for _, keyword := range keywords {
		if !slices.Contains(segments, keyword) {
			return false
		}
	}
	return true
}

func scoreExactKeywords(segments, exactKeywords []string) similarity {
	if len(segments) != 1 || len(exactKeywords) == 0 {
		return similarity{}
	}
	segment := segments[0]
	if slices.Contains(exactKeywords, segment) {
		return similarity{
			score:   1 * exactKeywordsBoost,
			matches: []string{segment},
		}
	}
	return similarity{}
}
