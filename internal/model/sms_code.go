package model

import "time"

type SmsCode struct {
	BaseModel
	Phone     string    `json:"phone" gorm:"size:20;index"`
	Code      string    `json:"-" gorm:"size:10"`
	Purpose   string    `json:"purpose" gorm:"size:20"`
	ExpiredAt time.Time `json:"expired_at"`
	Used      bool      `json:"used" gorm:"default:false"`
}
