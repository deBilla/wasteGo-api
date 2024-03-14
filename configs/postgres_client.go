package configs

import (
	"billacode/wasteGo/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=ziggy.db.elephantsql.com user=mkszuyvs password=RbIoPPH9H-kG4vPU3c6Sj-eBY0vrYuCg dbname=mkszuyvs"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&models.WasteItem{})
	if err != nil {
		return
	}

	DB = database
}
