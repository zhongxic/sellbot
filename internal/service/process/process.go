package process

import "time"

// Domain categories
const (
	DomainCategoryMainProcess  = "main_process"
	DomainCategoryBusinessQA   = "business_qa"
	DomainCategoryCommonDialog = "common_dialog"
	DomainCategorySilence      = "silence"
)

// Domain types
const (
	DomainTypeStart               = "start"
	DomainTypeNormal              = "normal"
	DomainTypeEnd                 = "end"
	DomainTypeAgent               = "agent"
	DomainTypeDialogEndFail       = "end_fail"
	DomainTypeDialogEndExceed     = "end_exceed"
	DomainTypeDialogClarification = "clarification"
)

// Branch semantics
const (
	BranchSemanticPositive = "positive"
	BranchSemanticNegative = "negative"
	BranchSemanticSpecial  = "special"
)

// Specific domain and branch names
const (
	DomainNameRepeat = "repeat"
	BranchNameEnter  = "enter"
)

// Interruption types
const (
	InterruptionTypeNone = iota
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
	Type            string            `json:"type"`
	Category        string            `json:"category"`
	Branches        map[string]Branch `json:"branches"`
	MatchOrders     []MatchPath       `json:"matchOrders"`
	IgnoreConfig    IgnoreConfig      `json:"ignoreConfig"`
	MissMatchConfig MissMatchConfig   `json:"missMatchConfig"`
}

type Branch struct {
	Name             string     `json:"name"`
	Keywords         Keywords   `json:"keywords"`
	Responses        []Response `json:"responses"`
	EnableExceedJump bool       `json:"enableExceedJump"`
	Next             string     `json:"next"`
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
}

type IntentionCondition struct {
	EnableCondition bool     `json:"enableCondition"`
	DomainName      string   `json:"domainName"`
	Keywords        Keywords `json:"keywords"`
}

type Options struct {
	MaxRounds int `json:"maxRounds"`
}

type Variable struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}
