package bot

type PrologueDTO struct {
	ProcessId string
	Variables map[string]string
	Test      bool
}

type ChatDTO struct {
	SessionId    string
	Sentence     string
	Silence      bool
	Interruption int
}

type InteractiveRespond struct {
	SessionId  string
	Hits       HitsDTO
	Answer     AnswerDTO
	Intentions []IntentionDTO
}

type HitsDTO struct {
	Sentence string
	Segments []string
	HitPaths []HitPathDTO
}

type HitPathDTO struct {
	Domain       string
	Branch       string
	MatchedWords []string
}

type AnswerDTO struct {
	Text  string
	Audio string
}

type IntentionDTO struct {
	Code        string
	DisplayName string
	Reason      string
}
