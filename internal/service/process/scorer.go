package process

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/zhongxic/sellbot/internal/traceid"
)

const (
	simpleKeywordBoost       = 1
	combinationKeywordsBoost = 2
	exactKeywordsBoost       = 100
)

type Similarity struct {
	Score   int
	Matches []string
}

func (s Similarity) IsMatched() bool {
	return s.Score > 0
}

func (s Similarity) IsBetterThan(other Similarity) bool {
	return s.Score > other.Score
}

func (s Similarity) String() string {
	return fmt.Sprintf("{score: %v, matches: %v}", s.Score, s.Matches)
}

func Score(ctx context.Context, text string, segments []string, Keywords Keywords) Similarity {
	simpleKeywordsSim := scoreSimpleKeywords(segments, Keywords.Simple)
	combinationKeywordsSim := scoreCombinationKeywords(segments, Keywords.Combination)
	exactKeywordsSim := scoreExactKeywords(segments, Keywords.Exact)
	slog.Debug(fmt.Sprintf("text [%v] segments [%v]: "+
		"simple keywords similarity [%v], combination keywords similarity [%v], exact keywords similarity [%v]",
		text, segments, simpleKeywordsSim, combinationKeywordsSim, exactKeywordsSim),
		slog.Any("traceId", ctx.Value(traceid.TraceId{})))
	maxSim := Similarity{}
	if simpleKeywordsSim.IsBetterThan(maxSim) {
		maxSim = simpleKeywordsSim
	}
	if combinationKeywordsSim.IsBetterThan(maxSim) {
		maxSim = combinationKeywordsSim
	}
	if exactKeywordsSim.IsBetterThan(maxSim) {
		maxSim = exactKeywordsSim
	}
	slog.Debug(fmt.Sprintf("text [%v] segments [%v]: best matched keywords similarity is [%v]",
		text, segments, maxSim),
		slog.Any("traceId", ctx.Value(traceid.TraceId{})))
	return maxSim
}

func scoreSimpleKeywords(segments, simpleKeywords []string) Similarity {
	if len(segments) == 0 || len(simpleKeywords) == 0 {
		return Similarity{}
	}
	var matches []string
	for _, segment := range segments {
		if slices.Contains(simpleKeywords, segment) {
			matches = append(matches, segment)
		}
	}
	return Similarity{
		Score:   len(matches) * simpleKeywordBoost,
		Matches: matches,
	}
}

func scoreCombinationKeywords(segments []string, combinationKeywords [][]string) Similarity {
	maxSim := Similarity{}
	if len(segments) == 0 || len(combinationKeywords) == 0 {
		return maxSim
	}
	for _, combinations := range combinationKeywords {
		if keywordsMatchAll(segments, combinations) {
			sim := Similarity{
				Score:   len(combinations) * combinationKeywordsBoost,
				Matches: combinations,
			}
			if sim.IsBetterThan(maxSim) {
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

func scoreExactKeywords(segments, exactKeywords []string) Similarity {
	if len(segments) != 1 || len(exactKeywords) == 0 {
		return Similarity{}
	}
	segment := segments[0]
	if slices.Contains(exactKeywords, segment) {
		return Similarity{
			Score:   1 * exactKeywordsBoost,
			Matches: []string{segment},
		}
	}
	return Similarity{}
}
