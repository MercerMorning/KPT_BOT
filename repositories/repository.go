package repositories

import (
	"KPT_BOT/database"
	"gorm.io/gorm"
)

var DB *gorm.DB = database.Init()
