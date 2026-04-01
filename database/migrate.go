package database

import "github.com/devjoemedia/go-ticketing-api/models"

func Migrate() {
	DB.AutoMigrate(
		&models.User{},
		&models.Todo{},
		&models.Ticket{},
		&models.RefreshToken{},
	)
}
