package models

type InvitationCategory struct {
	BaseModel
	IsActive bool   `gorm:"default:true;index"`
	Template string `gorm:"size:255;not null"` // HTML template path
	Name     string `gorm:"size:255;not null;index"`
	Icon     string `gorm:"size:50;not null"`  // Font Awesome icon class
	
	// Relationships
	Invitations []Invitation `gorm:"foreignKey:CategoryID"`
}

// TableName returns the table name for the InvitationCategory model
func (InvitationCategory) TableName() string {
	return "invitation_categories"
}
