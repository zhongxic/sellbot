package bot

import "github.com/zhongxic/sellbot/internal/service/bot"

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
}

type IntentionResponse struct {
	Code        string `json:"code"`
	DisplayName string `json:"displayName"`
	Reason      string `json:"reason"`
}

func convertPrologueRequestToPrologueDTO(request *PrologueRequest) *bot.PrologueDTO {
	return &bot.PrologueDTO{
		ProcessId: request.ProcessId,
		Variables: request.Variables,
		Test:      request.Test,
	}
}

func convertChatRequestToChatDTO(request *ChatRequest) *bot.ChatDTO {
	return &bot.ChatDTO{
		SessionId:    request.SessionId,
		Sentence:     request.Sentence,
		Silence:      request.Silence,
		Interruption: request.Interruption,
	}
}

func convertInteractiveRespondToInteractiveResponse(respond *bot.InteractiveRespond) *InteractiveResponse {
	return &InteractiveResponse{
		SessionId:  respond.SessionId,
		Hits:       convertHitsDTOToHitsResponse(respond.Hits),
		Answer:     convertAnswerDTOToAnswerResponse(respond.Answer),
		Intentions: convertIntentionDTOListToIntentionResponseList(respond.Intentions),
	}
}

func convertHitsDTOToHitsResponse(hitsDTO bot.HitsDTO) HitsResponse {
	hitsResponse := HitsResponse{}
	hitsResponse.Sentence = hitsDTO.Sentence
	if len(hitsDTO.Segments) == 0 {
		hitsResponse.Segments = make([]string, 0)
	} else {
		hitsResponse.Segments = hitsDTO.Segments
	}
	hitsResponse.HitPaths = convertHitPathDTOListToHitPathResponseList(hitsDTO.HitPaths)
	return hitsResponse
}

func convertHitPathDTOListToHitPathResponseList(hitPathDTOList []bot.HitPathDTO) []HitPathResponse {
	if len(hitPathDTOList) == 0 {
		return make([]HitPathResponse, 0)
	}
	hitPathResponseList := make([]HitPathResponse, len(hitPathDTOList))
	for i, hitPathDTO := range hitPathDTOList {
		hitPathResponseList[i] = convertHitPathDTOToHitPathResponse(hitPathDTO)
	}
	return hitPathResponseList
}

func convertHitPathDTOToHitPathResponse(hitPathDTO bot.HitPathDTO) HitPathResponse {
	hitPathResponse := HitPathResponse{}
	hitPathResponse.Domain = hitPathDTO.Domain
	hitPathResponse.Branch = hitPathDTO.Branch
	if len(hitPathDTO.MatchedWords) == 0 {
		hitPathResponse.MatchedWords = make([]string, 0)
	} else {
		hitPathResponse.MatchedWords = hitPathDTO.MatchedWords
	}
	return hitPathResponse
}

func convertAnswerDTOToAnswerResponse(answer bot.AnswerDTO) AnswerResponse {
	return AnswerResponse{
		Text:  answer.Text,
		Audio: answer.Audio,
	}
}

func convertIntentionDTOListToIntentionResponseList(intentionDTOList []bot.IntentionDTO) []IntentionResponse {
	if len(intentionDTOList) == 0 {
		return make([]IntentionResponse, 0)
	}
	intentionResponseList := make([]IntentionResponse, len(intentionDTOList))
	for i, intentionDTO := range intentionDTOList {
		intentionResponseList[i] = convertIntentionListToIntentionResponse(intentionDTO)
	}
	return intentionResponseList
}

func convertIntentionListToIntentionResponse(intentionDTO bot.IntentionDTO) IntentionResponse {
	return IntentionResponse{
		Code:        intentionDTO.Code,
		DisplayName: intentionDTO.DisplayName,
		Reason:      intentionDTO.Reason,
	}
}
