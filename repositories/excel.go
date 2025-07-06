package repositories

import (
	"KPT_BOT/models"
)

func SetExcel(chatId int64, token string, sheet string) error {
	task := models.Sheet{
		ChatId:  chatId,
		Code:    token,
		SheetId: sheet,
		Range:   60,
	}
	if result := DB.Create(&task); result.Error != nil {
		return result.Error
	}
	return nil
}

func ChangeTable(chatId int64, sheet string) error {
	task := models.Sheet{
		ChatId: chatId,
	}
	DB.Model(&task).Update("SheetId", sheet)
	return nil
}

func GetExcel(chatId int64) (models.Sheet, error) {
	var excel models.Sheet
	if result := DB.Where("chat_id = ?", chatId).Find(&excel); result.Error != nil {
		return excel, result.Error
	}
	return excel, nil
}

func GetAllIds() ([]int64, error) {
	var chatIds []int64
	err := DB.Model(&models.Sheet{}).Pluck("chat_id", &chatIds).Error
	if err != nil {
		return nil, err
	}
	return chatIds, nil
}
