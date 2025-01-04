package bot

import "github.com/zhongxic/sellbot/internal/service/bot"

func convertPrologueRequestToPrologueDTO(request *PrologueRequest) *bot.PrologueDTO {
	return &bot.PrologueDTO{
		ProcessId: request.ProcessId,
		Variables: request.Variables,
		Test:      request.Test,
	}
}

func convertConnectRequestToConnectDTO(request *ConnectRequest) *bot.ConnectDTO {
	return &bot.ConnectDTO{
		SessionId: request.SessionId,
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
		Ended: answer.Ended,
		Agent: answer.Agent,
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

func convertConnectRespondToResponse(respond *bot.ConnectRespond) *ConnectResponse {
	return &ConnectResponse{
		SessionId:  respond.SessionId,
		AnswerTime: respond.AnswerTime,
	}
}
