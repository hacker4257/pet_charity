package model

import "time"

type Organization struct {
	BaseModel
	UserID        uint       `json:"user_id" gorm:"uniqueIndex"`
	Name          string     `json:"name" gorm:"size:100"`
	LicenseNo     string     `json:"license_no" gorm:"size:100"`
	Description   string     `json:"description" gorm:"type:text"`
	Address       string     `json:"address" gorm:"size:500"`
	ContactPhone  string     `json:"contact_phone" gorm:"size:20"`
	OpeningHours  string     `json:"opening_hours" gorm:"size:100"`
	AcceptSpecies string     `json:"accept_species" gorm:"size:200"`
	Capacity      int        `json:"capacity" gorm:"default:0"`
	Longitude     float64    `json:"longitude" gorm:"type:decimal(10,7)"`
	Latitude      float64    `json:"latitude" gorm:"type:decimal(10,7)"`
	Status        string     `json:"status" gorm:"size:20;default:pending;index"`
	VerifiedAt    *time.Time `json:"verified_at"`

	//关联
	User User `json:"user" gorm:"foreignKey:UserID"`
}
