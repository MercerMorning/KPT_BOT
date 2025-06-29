package models

type Sheet struct {
	ChatId  int64  `gorm:"chat_id"`
	Code    string `gorm:"code"`
	SheetId string `gorm:"sheet"`
}
