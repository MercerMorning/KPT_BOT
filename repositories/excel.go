package repositories

import (
	"KPT_BOT/models"
)

func SetExcel(chatId int64, token string, sheet string) error {
	task := models.Sheet{
		ChatId:  chatId,
		Code:    token,
		SheetId: sheet,
	}
	if result := DB.Create(&task); result.Error != nil {
		return result.Error
	}
	return nil
}

func GetExcel(chatId int64) (models.Sheet, error) {
	var excel models.Sheet
	if result := DB.Where("chat_id = ?", chatId).Find(&excel); result.Error != nil {
		return excel, result.Error
	}
	return excel, nil
}
