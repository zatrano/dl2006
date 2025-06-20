package migrations

import (
	"davet.link/models"
	"gorm.io/gorm"
)

func CreateInvitationsTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.Invitation{})
}
