package bot

import "github.com/zhongxic/sellbot/internal/service/bot"

type Controller struct {
	botService bot.Service
}

func NewController(botService bot.Service) *Controller {
	return &Controller{botService: botService}
}
