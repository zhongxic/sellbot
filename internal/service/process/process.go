package process

import (
	"context"
	"strings"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

// DomainCategory is category of Domain.
type DomainCategory string

const (
	DomainCategoryMainProcess  DomainCategory = "main_process"
	DomainCategoryBusinessQA   DomainCategory = "business_qa"
	DomainCategoryCommonDialog DomainCategory = "common_dialog"
	DomainCategorySilence      DomainCategory = "silence"
)

// DomainType is type of Domain.
type DomainType string

func (t DomainType) IsEnded() bool {
	s := string(t)
	if t == DomainTypeEnd || strings.HasPrefix(s, "end_") {
		return true
	}
	return false
}

const (
	DomainTypeStart               DomainType = "start"
	DomainTypeNormal              DomainType = "normal"
	DomainTypeEnd                 DomainType = "end"
	DomainTypeAgent               DomainType = "agent"
	DomainTypeDialogConfused      DomainType = "confused"
	DomainTypeDialogRefused       DomainType = "refused"
	DomainTypeDialogMissMatch     DomainType = "miss_match"
	DomainTypeDialogEndFail       DomainType = "end_fail"
	DomainTypeDialogEndBusy       DomainType = "end_busy"
	DomainTypeDialogEndExceed     DomainType = "end_exceed"
	DomainTypeDialogEndMissMatch  DomainType = "end_miss_match"
	DomainTypeDialogEndException  DomainType = "end_exception"
	DomainTypeDialogCompliant     DomainType = "compliant"
	DomainTypeDialogPhoneFilter   DomainType = "phone_filter"
	DomainTypeDialogClarification DomainType = "clarification"
)

// BranchSemantic is semantics of Branch.
type BranchSemantic string

const (
	BranchSemanticPositive BranchSemantic = "positive"
	BranchSemanticNegative BranchSemantic = "negative"
	BranchSemanticSpecial  BranchSemantic = "special"
)

// Specific Domain and Branch names.
const (
	DomainNameRepeat = "repeat"
	BranchNameEnter  = "enter"
)

// Interruption is type of interruption.
type Interruption int

func (i Interruption) Value() int {
	return int(i)
}

const (
	InterruptionTypeNone Interruption = iota
	InterruptionTypeForce
	InterruptionTypeQA
	InterruptionTypeClarification
	InterruptionTypePrologue
)

type Process struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	Version    string            `json:"version"`
	Domains    map[string]Domain `json:"domains"`
	Intentions Intentions        `json:"intentions"`
	Options    Options           `json:"options"`
	Variables  []Variable        `json:"variables"`

	lastModified time.Time
}

type Domain struct {
	Name            string            `json:"name"`
	Type            DomainType        `json:"type"`
	Category        DomainCategory    `json:"category"`
	Order           int               `json:"order"`
	Branches        map[string]Branch `json:"branches"`
	MatchOrders     []MatchPath       `json:"matchOrders"`
	IgnoreConfig    IgnoreConfig      `json:"ignoreConfig"`
	MissMatchConfig MissMatchConfig   `json:"missMatchConfig"`
}

type Branch struct {
	Name             string         `json:"name"`
	Order            int            `json:"order"`
	Semantic         BranchSemantic `json:"semantic"`
	Keywords         Keywords       `json:"keywords"`
	Responses        []Response     `json:"responses"`
	EnableExceedJump bool           `json:"enableExceedJump"`
	Next             string         `json:"next"`
}

type Keywords struct {
	Simple      []string   `json:"simple"`
	Combination [][]string `json:"combination"`
	Exact       []string   `json:"exact"`
}

type Response struct {
	Text           string `json:"text"`
	Audio          string `json:"audio"`
	EnableAutoJump bool   `json:"enableAutoJump"`
	Next           string `json:"next"`
}

type MatchPath struct {
	DomainName string `json:"domainName"`
	BranchName string `json:"branchName"`
}

type IgnoreConfig struct {
	IgnoreAny              bool     `json:"ignoreAny"`
	IgnoreAnyExceptRefuse  bool     `json:"ignoreAnyExceptRefuse"`
	IgnoreAnyExceptDomains []string `json:"ignoreAnyExceptDomains"`
}

type MissMatchConfig struct {
	LongTextMissMatchJumpTo  string `json:"longTextMissMatchJumpTo"`
	ShortTextMissMatchJumpTo string `json:"shortTextMissMatchJumpTo"`
}

type Intentions struct {
	DefaultIntention string          `json:"defaultIntention"`
	IntentionRules   []IntentionRule `json:"intentionRules"`
}

type IntentionRule struct {
	Code               string             `json:"code"`
	Expression         string             `json:"expression"`
	DisplayName        string             `json:"displayName"`
	Reason             string             `json:"reason"`
	IntentionCondition IntentionCondition `json:"intentionCondition"`

	// program is an expression who compiled from Expression.
	program *vm.Program
}

func (i IntentionRule) IsHit(ctx context.Context, env IntentionAnalyzeEnv) (bool, error) {
	expressionHit := true
	intentionConditionHit := true
	// check if expression matched.
	if i.program != nil {
		value, err := expr.Run(i.program, env)
		if err != nil {
			return false, err
		}
		if v, ok := value.(bool); ok {
			expressionHit = v
		} else {
			expressionHit = false
		}
	}
	// check if intention condition matched.
	if i.IntentionCondition.Enabled {
		domainMatched := env.Status.PreviousMainProcessDomain == i.IntentionCondition.DomainName
		similarity := Score(ctx, env.Sentence, env.Segments, i.IntentionCondition.Keywords)
		intentionConditionHit = domainMatched && similarity.IsMatched()
	}
	return expressionHit && intentionConditionHit, nil
}

type IntentionCondition struct {
	Enabled    bool     `json:"enabled"`
	DomainName string   `json:"domainName"`
	Keywords   Keywords `json:"keywords"`
}

type Options struct {
	MaxRounds              int    `json:"maxRounds"`
	ForceInterruptedJumpTo string `json:"forceInterruptedJumpTo"`
}

type Variable struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}
