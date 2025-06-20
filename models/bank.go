package models

type Bank struct {
	BaseModel
	IsActive bool   `gorm:"default:true;index"`
	Name     string `gorm:"size:255;not null;index"`
}

// TableName returns the table name for the Bank model
func (Bank) TableName() string {
	return "banks"
}
