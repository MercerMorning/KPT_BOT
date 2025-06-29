package services

import (
	"KPT_BOT/clients"
	"KPT_BOT/repositories"
	"KPT_BOT/session"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/oauth2"
	"time"
)

func Start(bot *tgbotapi.BotAPI, update tgbotapi.Update, usrSession *session.Session) {
	usrSession.Stage = session.RequestGoogleSheetsApiKey
	text := "Привет. Это дневник КПТ. Сейчас мы перейдем к настройке. Введите uuid таблицы в google sheets"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func Write(bot *tgbotapi.BotAPI, update tgbotapi.Update, usrSession *session.Session) {
	usrSession.Stage = session.WritingDiarySituation

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Опишите ситуацию")
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func GetCodeFromWeb(bot *tgbotapi.BotAPI, update tgbotapi.Update, usrSession *session.Session) {
	usrSession.SheetId = update.Message.Text
	gs := clients.SheetsClient{
		Update: update,
		Bot:    bot,
	}
	gs.RequestCode()
	usrSession.Stage = session.ReceiveGoogleSheetsApiKey
}

func InitTable(bot *tgbotapi.BotAPI, update tgbotapi.Update, usrSession *session.Session) {
	gs := clients.SheetsClient{
		Update: update,
		Bot:    bot,
	}
	tok := gs.GetToken()
	token, err := json.Marshal(tok)
	if err != nil {
		fmt.Printf("Unable to marshall token: %v", err)
	}

	err = repositories.SetExcel(update.Message.Chat.ID, string(token), usrSession.SheetId)

	if err != nil {
		fmt.Printf("Unable to set excel: %v", err)
	}

	gs.InitTable(tok, usrSession.SheetId)

	usrSession.Stage = session.WritingDiarySituation

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Опишите ситуацию")
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func Append(bot *tgbotapi.BotAPI, update tgbotapi.Update, usrSession *session.Session) {
	gs := clients.SheetsClient{
		Update: update,
		Bot:    bot,
	}
	token := usrSession.Code
	tok := oauth2.Token{}
	err := json.Unmarshal([]byte(token), &tok)
	if err != nil {
		fmt.Printf("Unable to unmarshall token: %v", err)
	}

	err = repositories.SetExcel(update.Message.Chat.ID, string(token), usrSession.SheetId)

	if err != nil {
		fmt.Printf("Unable to set excel: %v", err)
	}

	data := []string{
		usrSession.Diary.Situation,
		usrSession.Diary.Thought,
		usrSession.Diary.Emotion,
		usrSession.Diary.Feeling,
		usrSession.Diary.Action,
		time.Now().Format(time.DateOnly),
	}

	gs.Append(&tok, usrSession.SheetId, data)

	usrSession.Stage = session.WritingDiarySituation

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Опишите ситуацию")
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}
