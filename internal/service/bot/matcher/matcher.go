package matcher

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/zhongxic/sellbot/internal/service/process"
	"github.com/zhongxic/sellbot/internal/traceid"
)

// DefaultChainedMatcher default matcher chain
var DefaultChainedMatcher = &ChainedMatcher{
	matchers: []Matcher{
		&OutOfMaxRoundsMatcher{},
		&ForceInterruptionMatcher{},
		&ClarificationInterruptionMatcher{},
		&SilenceMatcher{},
		&PreIgnoreMatcher{},
		&TextMatcher{},
		&PostIgnoreMatcher{},
		&MissMatchMatcher{},
	},
}

type Matcher interface {
	// Match find matched branches in process use this matchContext
	// and abort match if return false
	Match(ctx context.Context, matchContext *Context) (abort bool, err error)
}

type OutOfMaxRoundsMatcher struct {
}

func (matcher *OutOfMaxRoundsMatcher) Match(ctx context.Context, matchContext *Context) (bool, error) {
	if matchContext.Session.ConversationCount >= matchContext.Process.Options.MaxRounds {
		slog.Info(fmt.Sprintf("sessionId [%v]: OutOfMaxRoundsMatcher detect conversation count out of max rounds [%v]",
			matchContext.Session.Id, matchContext.Process.Options.MaxRounds),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		processHelper := process.NewHelper(matchContext.Process)
		domain, err := processHelper.GetCommonDialog(process.DomainTypeDialogEndExceed)
		if err != nil {
			return true, fmt.Errorf("OutOfMaxRoundsMatcher get common dialog [%v] failed: %w", process.DomainTypeDialogEndExceed, err)
		}
		matchedPath := MatchedPath{Domain: domain.Name, Branch: process.BranchNameEnter}
		slog.Info(fmt.Sprintf("sessionId [%v]: OutOfMaxRoundsMatcher matched domain [%v] branch [%v]",
			matchContext.Session.Id, matchedPath.Domain, matchedPath.Branch),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		matchContext.AddMatchedPath(matchedPath)
		return true, nil
	}
	return false, nil
}

type ForceInterruptionMatcher struct {
}

func (matcher *ForceInterruptionMatcher) Match(ctx context.Context, matchContext *Context) (bool, error) {
	if matchContext.Interruption == process.InterruptionTypeForce.Value() {
		slog.Info(fmt.Sprintf("sessionId [%v]: ForceInterruptionMatcher detect force interruption", matchContext.Session.Id),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		if matchContext.Process.Options.ForceInterruptedJumpTo == "" {
			slog.Info(fmt.Sprintf("sessionId [%v]: ForceInterruptionMatcher not handle interruption because jump to domain is not specified",
				matchContext.Session.Id),
				slog.Any("traceId", ctx.Value(traceid.TraceId{})))
			return false, nil
		}
		matchedPath := MatchedPath{Domain: matchContext.Process.Options.ForceInterruptedJumpTo, Branch: process.BranchNameEnter}
		slog.Info(fmt.Sprintf("sessionId [%v]: ForceInterruptionMatcher matched domain [%v] branch [%v]",
			matchContext.Session.Id, matchedPath.Domain, matchedPath.Branch),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		matchContext.AddMatchedPath(matchedPath)
		return true, nil
	}
	return false, nil
}

type ClarificationInterruptionMatcher struct {
}

func (matcher *ClarificationInterruptionMatcher) Match(ctx context.Context, matchContext *Context) (bool, error) {
	if matchContext.Interruption == process.InterruptionTypeClarification.Value() {
		slog.Info(fmt.Sprintf("sessionId [%v]: ClarificationInterruptionMatcher detect clarification interruption",
			matchContext.Session.Id),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		processHelper := process.NewHelper(matchContext.Process)
		domain, err := processHelper.GetCommonDialog(process.DomainTypeDialogClarification)
		if err != nil {
			return true, fmt.Errorf("ClarificationInterruptionMatcher get common dialog [%v] failed: %w",
				process.DomainTypeDialogClarification, err)
		}
		matchedPath := MatchedPath{Domain: domain.Name, Branch: process.BranchNameEnter}
		slog.Info(fmt.Sprintf("sessionId [%v]: ClarificationInterruptionMatcher matched domain [%v] branch [%v]",
			matchContext.Session.Id, matchedPath.Domain, matchedPath.Branch),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		matchContext.AddMatchedPath(matchedPath)
		return true, nil
	}
	return false, nil
}

type SilenceMatcher struct {
}

func (matcher *SilenceMatcher) Match(ctx context.Context, matchContext *Context) (bool, error) {
	if matchContext.Silence {
		slog.Info(fmt.Sprintf("sessionId [%v]: SilenceMatcher detect silence", matchContext.Session.Id),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		processHelper := process.NewHelper(matchContext.Process)
		domain, err := processHelper.GetSilenceDomain()
		if err != nil {
			return true, fmt.Errorf("SilenceMatcher get silence domain failed: %w", err)
		}
		matchedPath := MatchedPath{Domain: domain.Name, Branch: process.BranchNameEnter}
		slog.Info(fmt.Sprintf("sessionId [%v]: SilenceMatcher matched domain [%v] branch [%v]",
			matchContext.Session.Id, matchedPath.Domain, matchedPath.Branch),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		matchContext.AddMatchedPath(matchedPath)
		return true, nil
	}
	return false, nil
}

type PreIgnoreMatcher struct {
}

func (matcher *PreIgnoreMatcher) Match(ctx context.Context, matchContext *Context) (bool, error) {
	processHelper := process.NewHelper(matchContext.Process)
	domain, err := processHelper.GetDomain(matchContext.Session.CurrentDomain)
	if err != nil {
		return true, fmt.Errorf("PreIgnoreMatcher get current [%v] domain failed: %w",
			matchContext.Session.CurrentDomain, err)
	}
	if domain.IgnoreConfig.IgnoreAny {
		slog.Info(fmt.Sprintf("sessionId [%v]: PreIgnoreMatcher detect current domain [%v] turned on ignore any",
			matchContext.Session.Id, domain.Name),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		branch, err := processHelper.GetDomainPositiveBranch(domain.Name)
		if err != nil {
			return true, fmt.Errorf("PreIgnoreMatcher get current domain [%v] positive branch failed: %w",
				matchContext.Session.CurrentDomain, err)
		}
		matchedPath := MatchedPath{Domain: domain.Name, Branch: branch.Name}
		slog.Info(fmt.Sprintf("sessionId [%v]: PreIgnoreMatcher matched domain [%v] branch [%v]",
			matchContext.Session.Id, matchedPath.Domain, matchedPath.Branch),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		matchContext.AddMatchedPath(matchedPath)
		return true, nil
	}
	return false, nil
}

type TextMatcher struct {
}

func (matcher *TextMatcher) Match(ctx context.Context, matchContext *Context) (bool, error) {
	processHelper := process.NewHelper(matchContext.Process)
	matchPaths, err := processHelper.GetMergeOrderedMatchPaths(matchContext.Session.CurrentMainProcessDomain)
	if err != nil {
		return true, fmt.Errorf("TextMatcher get domain [%v] merge ordered match paths failed: %w",
			matchContext.Session.CurrentMainProcessDomain, err)
	}
	maxSimilarity := process.Similarity{}
	bestMatchedPath := process.MatchPath{}
	for _, matchPath := range matchPaths {
		branch, err := processHelper.GetBranch(matchPath.DomainName, matchPath.BranchName)
		if err != nil {
			return true, fmt.Errorf("TextMatcher get domain [%v] branch [%v] failed: %w",
				matchPath.DomainName, matchPath.BranchName, err)
		}
		similarity := process.Score(ctx, matchContext.Sentence, matchContext.Segments, branch.Keywords)
		slog.Debug(fmt.Sprintf("sessionId [%v]: TextMatcher current mainProcessDomain [%v] "+
			"detect domain [%v] branch [%v] similarity score [%v] isMatched [%v]",
			matchContext.Session.Id, matchContext.Session.CurrentMainProcessDomain,
			matchPath.DomainName, matchPath.BranchName, similarity.Score, similarity.IsMatched()),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))

		if similarity.IsMatched() && similarity.IsBetterThan(maxSimilarity) {
			maxSimilarity = similarity
			bestMatchedPath = matchPath
		}
	}
	if maxSimilarity.IsMatched() {
		slog.Info(fmt.Sprintf("sessionId [%v]: TextMatcher current mainProcessDomain [%v] "+
			"detect best matched domain [%v] branch [%v]",
			matchContext.Session.Id, matchContext.Session.CurrentMainProcessDomain,
			bestMatchedPath.DomainName, bestMatchedPath.BranchName),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		matchContext.AddMatchedPath(MatchedPath{
			Domain:       bestMatchedPath.DomainName,
			Branch:       bestMatchedPath.BranchName,
			MatchedWords: maxSimilarity.Matches,
		})
	} else {
		slog.Info(fmt.Sprintf("sessionId [%v]: TextMatcher current MainProcessDomain [%v] detect miss match",
			matchContext.Session.Id, matchContext.Session.CurrentMainProcessDomain),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		domain, err := processHelper.GetCommonDialog(process.DomainTypeDialogMissMatch)
		if err != nil {
			return true, fmt.Errorf("TextMatcher get common dialog [%v] failed: %w", process.DomainTypeDialogMissMatch, err)
		}
		matchedPath := MatchedPath{Domain: domain.Name, Branch: process.BranchNameEnter}
		slog.Info(fmt.Sprintf("sessionId [%v]: TextMatcher current MainProcessDomain [%v] matched domain [%v] branch [%v]",
			matchContext.Session.Id, matchContext.Session.CurrentMainProcessDomain, matchedPath.Domain, matchedPath.Branch),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		matchContext.AddMatchedPath(matchedPath)
	}
	return false, nil
}

type PostIgnoreMatcher struct {
}

func (matcher *PostIgnoreMatcher) Match(ctx context.Context, matchContext *Context) (bool, error) {
	lastMatchedPath, err := matchContext.GetLastMatchedPath()
	if err != nil {
		return true, fmt.Errorf("PostIgnoreMatcher get last matched path failed: %w", err)
	}
	processHelper := process.NewHelper(matchContext.Process)
	domain, err := processHelper.GetDomain(matchContext.Session.CurrentDomain)
	if err != nil {
		return true, fmt.Errorf("PostIgnoreMatcher get current domain [%v] failed: %w",
			matchContext.Session.CurrentDomain, err)
	}
	matchedDomain, err := processHelper.GetDomain(lastMatchedPath.Domain)
	if err != nil {
		return true, fmt.Errorf("PostIgnoreMatcher get last matched domain [%v] failed: %w", lastMatchedPath.Domain, err)
	}
	ignoreAnyExceptRefuse := domain.IgnoreConfig.IgnoreAnyExceptRefuse &&
		matchedDomain.Type != process.DomainTypeDialogRefused
	nextDomain := ""
	if matchedDomain.Category == process.DomainCategoryMainProcess {
		branch, err := processHelper.GetBranch(lastMatchedPath.Domain, lastMatchedPath.Branch)
		if err != nil {
			return true, fmt.Errorf("PostIgnoreMatcher get domain [%v] branch [%v] failed: %w",
				lastMatchedPath.Domain, lastMatchedPath.Branch, err)
		}
		nextDomain = branch.Next
	}
	if matchedDomain.Category == process.DomainCategoryBusinessQA {
		nextDomain = lastMatchedPath.Domain + "." + lastMatchedPath.Branch
	}
	if matchedDomain.Category == process.DomainCategoryCommonDialog {
		nextDomain = lastMatchedPath.Domain
	}
	ignoreAnyExceptDomains := false
	if len(domain.IgnoreConfig.IgnoreAnyExceptDomains) > 0 {
		ignoreAnyExceptDomains = !slices.Contains(domain.IgnoreConfig.IgnoreAnyExceptDomains, nextDomain)
	}
	shouldIgnore := ignoreAnyExceptRefuse || ignoreAnyExceptDomains
	if shouldIgnore {
		slog.Info(fmt.Sprintf("sessionId [%v]: PostIgnoreMatcher ignoreAnyExceptRefuse [%v] ignoreAnyExceptDomains [%v]",
			matchContext.Session.Id, ignoreAnyExceptRefuse, ignoreAnyExceptDomains),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		branch, err := processHelper.GetDomainPositiveBranch(matchContext.Session.CurrentDomain)
		if err != nil {
			return true, fmt.Errorf("PostIgnoreMatcher get current domain [%v] positive branch failed: %w",
				matchContext.Session.CurrentDomain, err)
		}
		matchedPath := MatchedPath{Domain: matchContext.Session.CurrentDomain, Branch: branch.Name}
		slog.Info(fmt.Sprintf("sessionId [%v]: PostIgnoreMatcher matched domain [%v] branch [%v]",
			matchContext.Session.Id, matchedPath.Domain, matchedPath.Branch),
			slog.Any("traceId", ctx.Value(traceid.TraceId{})))
		matchContext.AddMatchedPath(matchedPath)
	}
	return shouldIgnore, nil
}

type MissMatchMatcher struct {
}

func (matcher *MissMatchMatcher) Match(ctx context.Context, matchContext *Context) (bool, error) { // NOSONAR
	lastMatchedPath, err := matchContext.GetLastMatchedPath()
	if err != nil {
		return false, fmt.Errorf("MissMatchMatcher get last matched path failed: %w", err)
	}
	processHelper := process.NewHelper(matchContext.Process)
	matchedDomain, err := processHelper.GetDomain(lastMatchedPath.Domain)
	if err != nil {
		return true, fmt.Errorf("MissMatchMatcher get last matched domain [%v] failed: %w", lastMatchedPath.Domain, err)
	}
	domain, err := processHelper.GetDomain(matchContext.Session.CurrentMainProcessDomain)
	if err != nil {
		return true, fmt.Errorf("MissMatchMatcher get current mainProcessDomain [%v] failed: %w",
			matchContext.Session.CurrentDomain, err)
	}
	if matchedDomain.Type == process.DomainTypeDialogMissMatch {
		if matchContext.Session.MissMatchCount >= 3 {
			slog.Info("sessionId [%v]: MissMatchMatcher current main process domain [%v] detect miss match count exceed",
				matchContext.Session.Id, matchContext.Session.CurrentMainProcessDomain,
				slog.Any("traceId", ctx.Value(traceid.TraceId{})))
			dialog, err := processHelper.GetCommonDialog(process.DomainTypeDialogEndMissMatch)
			if err != nil {
				return true, fmt.Errorf("MissMatchMatcher get end_miss_match common dialog failed: %w", err)
			}
			matchedPath := MatchedPath{Domain: dialog.Name, Branch: process.BranchNameEnter}
			slog.Info(fmt.Sprintf("sessionId [%v]: MissMatchMatcher matched domain [%v] branch [%v]",
				matchContext.Session.Id, matchedPath.Domain, matchedPath.Branch),
				slog.Any("traceId", ctx.Value(traceid.TraceId{})))
			matchContext.AddMatchedPath(matchedPath)
			return true, nil
		}
		jumpTo := ""
		shortTextMissMatchJump := len([]rune(matchContext.Sentence)) < 4 && domain.MissMatchConfig.ShortTextMissMatchJumpTo != ""
		if shortTextMissMatchJump {
			slog.Info(fmt.Sprintf("sessionId [%v]: MissMatchMatcher current main process domain [%v] detect miss match [short text] jump to [%v]",
				matchContext.Session.Id, matchContext.Session.CurrentMainProcessDomain, domain.MissMatchConfig.ShortTextMissMatchJumpTo),
				slog.Any("traceId", ctx.Value(traceid.TraceId{})))
			jumpTo = domain.MissMatchConfig.ShortTextMissMatchJumpTo
		}
		longTextMissMatchJump := len([]rune(matchContext.Sentence)) >= 4 && domain.MissMatchConfig.LongTextMissMatchJumpTo != ""
		if longTextMissMatchJump {
			slog.Info(fmt.Sprintf("sessionId [%v]: MissMatchMatcher current main process domain [%v] detect miss match [long text] jump to [%v]",
				matchContext.Session.Id, matchContext.Session.CurrentMainProcessDomain, domain.MissMatchConfig.LongTextMissMatchJumpTo),
				slog.Any("traceId", ctx.Value(traceid.TraceId{})))
			jumpTo = domain.MissMatchConfig.LongTextMissMatchJumpTo
		}
		if jumpTo != "" {
			matchedPath := MatchedPath{Domain: jumpTo, Branch: process.BranchNameEnter}
			slog.Info(fmt.Sprintf("sessionId [%v]: MissMatchMatcher matched domain [%v] branch [%v]",
				matchContext.Session.Id, matchedPath.Domain, matchedPath.Branch),
				slog.Any("traceId", ctx.Value(traceid.TraceId{})))
			matchContext.AddMatchedPath(matchedPath)
		}
	}
	return true, nil
}

type ChainedMatcher struct {
	matchers []Matcher
}

func (c *ChainedMatcher) Match(ctx context.Context, matchContext *Context) (bool, error) {
	if len(c.matchers) == 0 {
		return true, fmt.Errorf("no matchers provide in chained matcher")
	}
	for _, matcher := range c.matchers {
		abort, err := matcher.Match(ctx, matchContext)
		if err != nil {
			return true, fmt.Errorf("chained matcher match failed: %w", err)
		}
		if abort {
			return true, nil
		}
	}
	return true, nil
}
