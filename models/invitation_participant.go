package models

type InvitationParticipant struct {
	BaseModel
	Title        string    `gorm:"size:255;not null"`
	PhoneNumber  string    `gorm:"size:20;not null"`
	GuestCount   int       `gorm:"not null;default:1"`
	InvitationID uint      `gorm:"index;not null"` // Foreign key for many-to-one relationship
	Invitation   Invitation
}

// TableName returns the table name for the InvitationParticipant model
func (InvitationParticipant) TableName() string {
	return "invitation_participants"
}
