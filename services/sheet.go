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

func RequestChangeTable(bot *tgbotapi.BotAPI, update tgbotapi.Update, usrSession *session.Session) {
	usrSession.Stage = session.ChangeTable

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите новую таблицу")
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func ChangeTable(update tgbotapi.Update, usrSession *session.Session) {
	usrSession.Stage = session.ChangeTable
	err := repositories.ChangeTable(update.Message.Chat.ID, update.Message.Text)
	if err != nil {
		panic(err)
	}
}

func GetCodeFromWeb(bot *tgbotapi.BotAPI, update tgbotapi.Update, usrSession *session.Session) {
	usrSession.SheetId = update.Message.Text
	gs := clients.SheetsClient{
		Update: update,
		Bot:    bot,
		ChatId: update.Message.Chat.ID,
	}
	gs.RequestCode()
	usrSession.Stage = session.ReceiveGoogleSheetsApiKey
}

func InitTable(code string, bot *tgbotapi.BotAPI, chatId int64, usrSession *session.Session) {
	gs := clients.SheetsClient{
		Code: code,
		//Update: update,
		Bot:    bot,
		ChatId: chatId,
	}
	tok := gs.GetToken()
	token, err := json.Marshal(tok)
	if err != nil {
		fmt.Printf("Unable to marshall token: %v", err)
	}
	fmt.Println("after get token")
	err = repositories.SetExcel(chatId, string(token), usrSession.SheetId)

	if err != nil {
		fmt.Printf("Unable to set excel: %v", err)
	}

	gs.InitTable(tok, usrSession.SheetId)

	usrSession.Stage = session.WritingDiarySituation

	msg := tgbotapi.NewMessage(chatId, "Опишите ситуацию")
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func Append(bot *tgbotapi.BotAPI, update tgbotapi.Update, usrSession *session.Session) {
	gs := clients.SheetsClient{
		Update: update,
		Bot:    bot,
		ChatId: update.Message.Chat.ID,
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

// Функция для отправки уведомления в один чат
func sendTelegramNotification(bot *tgbotapi.BotAPI, chatID int64, message string) error {
	// Здесь реализация отправки через Telegram Bot API
	// Например, используя github.com/go-telegram-bot-api/telegram-bot-api
	msg := tgbotapi.NewMessage(chatID, "Заполните дневник")
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
	//fmt.Printf("Отправляю уведомление в чат %d: %s\n", chatID, message)
	return nil
}

// Функция для отправки всем чатам
func notifyAllChats(bot *tgbotapi.BotAPI) {
	chatIds, err := repositories.GetAllIds()
	if err != nil {
		fmt.Printf("Ошибка получения chat_ids: %v\n", err)
		return
	}

	message := "Ваше регулярное уведомление" // Текст уведомления

	for _, chatID := range chatIds {
		if err := sendTelegramNotification(bot, chatID, message); err != nil {
			fmt.Printf("Ошибка отправки в чат %d: %v\n", chatID, err)
		}
		time.Sleep(500 * time.Millisecond) // Небольшая задержка между сообщениями
	}
}

func StartNotificationScheduler(bot *tgbotapi.BotAPI) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	// Первый запуск сразу
	notifyAllChats(bot)

	for range ticker.C {
		notifyAllChats(bot)
	}
}

//func RemindAboutDiary() {
//	chats :=
//}
