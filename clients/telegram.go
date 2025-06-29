package clients

import (
	"KPT_BOT/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Init() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(config.Config("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	return bot
}
