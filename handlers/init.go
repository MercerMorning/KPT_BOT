package handlers

import (
	"KPT_BOT/repositories"
	"KPT_BOT/services"
	"KPT_BOT/session"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Init(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	usrSessions := map[int64]*session.Session{}
	updates := bot.GetUpdatesChan(u)

	go services.StartNotificationScheduler(bot)

	for update := range updates {
		usrSession, ok := usrSessions[update.Message.Chat.ID]
		if !ok {
			excel, err := repositories.GetExcel(update.Message.Chat.ID)
			if err != nil {
				panic(err)
			}
			usrSession = &session.Session{
				session.Start,
				excel.Code,
				excel.SheetId,
				session.Diary{}}
			usrSessions[update.Message.Chat.ID] = usrSession
		}

		if usrSession.Code != "" && usrSession.Stage == session.Start {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Опишите ситуацию")
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
			usrSession.Stage = session.WritingDiarySituation
		}

		if update.CallbackQuery != nil {
			//Callbacks(bot, update)
		} else if update.Message.IsCommand() {
			Commands(bot, update, usrSession)
		} else {
			switch usrSession.Stage {
			case session.Start:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите команду")
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			case session.RequestGoogleSheetsApiKey:
				services.GetCodeFromWeb(bot, update, usrSession)
			case session.ReceiveGoogleSheetsApiKey:
				services.InitTable(bot, update, usrSession)
			case session.WritingDiarySituation:
				usrSession.Diary.Situation = update.Message.Text
				usrSession.Stage = session.WritingDiaryThought
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Опишите мысль")
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			case session.WritingDiaryThought:
				usrSession.Diary.Thought = update.Message.Text
				usrSession.Stage = session.WritingDiaryEmotion
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Опишите эмоцию")
				var numericKeyboard = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("Радость"),
						tgbotapi.NewKeyboardButton("Грусть"),
						tgbotapi.NewKeyboardButton("Гнев"),
					),
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("Страх"),
						tgbotapi.NewKeyboardButton("Удивление"),
						tgbotapi.NewKeyboardButton("Отвращение"),
					),
				)
				msg.ReplyMarkup = numericKeyboard
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			case session.WritingDiaryEmotion:
				usrSession.Diary.Emotion = update.Message.Text
				usrSession.Stage = session.WritingDiaryFeeling
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Опишите ощущение")
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			case session.WritingDiaryFeeling:
				usrSession.Diary.Feeling = update.Message.Text
				usrSession.Stage = session.WritingDiaryAction
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Действие")
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			case session.WritingDiaryAction:
				usrSession.Diary.Action = update.Message.Text
				services.Append(bot, update, usrSession)
			default:
				fmt.Println("Default")
			}
		}
	}
}
