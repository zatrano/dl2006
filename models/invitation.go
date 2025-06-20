package models

import "time"

type Invitation struct {
	BaseModel
	
	// Required fields
	InvitationKey  string    `gorm:"size:100;uniqueIndex;not null"` // Unique identifier for public access
	UserID         uint      `gorm:"index;not null"`                 // Foreign key for User relationship
	CategoryID     uint      `gorm:"index;not null"`                 // Foreign key for Category relationship
	Template       string    `gorm:"size:100;not null"`             // Template to use for rendering
	Type          string    `gorm:"size:50;not null;default:'basic'"` // Invitation type
	
	// Optional fields
	Title         string    `gorm:"size:255"`
	Image         string    `gorm:"size:255"`              // Cover image URL
	Description   string    `gorm:"type:text"`            // HTML content allowed
	Venue         string    `gorm:"size:255"`             // Event venue name
	Address       string    `gorm:"size:255"`             // Full address
	Location      string    `gorm:"size:255"`             // Google Maps link or coordinates
	Link          string    `gorm:"size:255"`             // External link if any
	Telephone     string    `gorm:"size:20"`              // Contact number
	Note          string    `gorm:"type:text"`            // Additional notes
	Date          time.Time `gorm:"index"`                // Event date
	Time          time.Time                               // Event time
	
	// Status fields
	IsConfirmed   bool      `gorm:"default:false;index"`  // Whether approved by admin
	IsParticipant bool      `gorm:"default:true"`         // Whether participation is allowed
	
	// Relationships
	User               *User                   `gorm:"foreignKey:UserID"`
	Category           *InvitationCategory     `gorm:"foreignKey:CategoryID"`
	InvitationDetail   *InvitationDetail      `gorm:"foreignKey:InvitationID;references:ID"`
	Participants       []InvitationParticipant `gorm:"foreignKey:InvitationID"`
}

// TableName returns the table name for the Invitation model
func (Invitation) TableName() string {
	return "invitations"
}
