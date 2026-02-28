package model

import "time"

type Adoption struct {
	BaseModel
	UserID          uint       `json:"user_id" gorm:"index"`
	PetID           uint       `json:"pet_id" gorm:"index"`
	Reason          string     `json:"reason" gorm:"type:text"`
	LivingCondition string     `json:"living_condition" gorm:"size:200"`
	Experience      string     `json:"experience" gorm:"type:text"`
	Status          string     `json:"status" gorm:"size:20;default:pending;index"`
	ReviewedBy      *uint      `json:"reviewed_by"`
	ReviewedAt      *time.Time `json:"reviewed_at"`
	RejectReason    string     `json:"reject_reason" gorm:"size:500"`

	//关联
	User User `json:"user" gorm:"foreignKey:UserID"`
	Pet  Pet  `json:"pet" gorm:"foreignKey:PetID"`
}
