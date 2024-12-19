package bot

import "github.com/zhongxic/sellbot/internal/service/bot"

type PrologueRequest struct {
	ProcessId string            `json:"processId"`
	Variables map[string]string `json:"variables"`
	Test      bool              `json:"test"`
}

type InteractiveResponse struct {
}

func convertPrologueRequestToPrologueDTO(request *PrologueRequest) *bot.PrologueDTO {
	return &bot.PrologueDTO{
		ProcessId: request.ProcessId,
		Variables: request.Variables,
		Test:      request.Test,
	}
}

func convertInteractiveRespondToInteractiveResponse(respond *bot.InteractiveRespond) *InteractiveResponse {
	return &InteractiveResponse{}
}
