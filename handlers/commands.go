package handlers

import (
	"KPT_BOT/services"
	"KPT_BOT/session"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Commands(bot *tgbotapi.BotAPI, update tgbotapi.Update, session *session.Session) {
	switch update.Message.Command() {
	case "start":
		services.Start(bot, update, session)
	case "write":
		services.Write(bot, update, session)
	}
}
