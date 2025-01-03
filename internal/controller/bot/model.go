package bot

type PrologueRequest struct {
	ProcessId string            `json:"processId"`
	Variables map[string]string `json:"variables"`
	Test      bool              `json:"test"`
}

type ChatRequest struct {
	SessionId    string `json:"sessionId"`
	Sentence     string `json:"sentence"`
	Silence      bool   `json:"silence"`
	Interruption int    `json:"interruption"`
}

type InteractiveResponse struct {
	SessionId  string              `json:"sessionId"`
	Hits       HitsResponse        `json:"hits"`
	Answer     AnswerResponse      `json:"answer"`
	Intentions []IntentionResponse `json:"intentions"`
}

type HitsResponse struct {
	Sentence string            `json:"sentence"`
	Segments []string          `json:"segments"`
	HitPaths []HitPathResponse `json:"hitPaths"`
}

type HitPathResponse struct {
	Domain       string   `json:"domain"`
	Branch       string   `json:"branch"`
	MatchedWords []string `json:"matchedWords"`
}

type AnswerResponse struct {
	Text  string `json:"text"`
	Audio string `json:"audio"`
	Ended bool   `json:"ended"`
	Agent bool   `json:"agent"`
}

type IntentionResponse struct {
	Code        string `json:"code"`
	DisplayName string `json:"displayName"`
	Reason      string `json:"reason"`
}
