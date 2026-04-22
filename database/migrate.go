package database

import "github.com/devjoemedia/scrumpilot-go-api/models"

func Migrate() {
	DB.AutoMigrate(
		&models.User{},
		&models.Todo{},
		&models.Ticket{},
		&models.RefreshToken{},
	)
}
