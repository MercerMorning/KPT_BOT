package database

import (
	"KPT_BOT/models"
	_ "github.com/lib/pq"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.Sheet{})
	if err != nil {
		panic(err)
	}

	return db
}
