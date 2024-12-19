package bot

import "github.com/zhongxic/sellbot/internal/service/bot"

type PrologueRequest struct {
}

type InteractiveResponse struct {
}

func convertPrologueRequestToPrologueDTO(request *PrologueRequest) *bot.PrologueDTO {
	return &bot.PrologueDTO{}
}

func convertInteractiveRespondToInteractiveResponse(respond *bot.InteractiveRespond) *InteractiveResponse {
	return &InteractiveResponse{}
}
