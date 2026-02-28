package model

import "time"

type Donation struct {
	BaseModel
	UserID        uint       `json:"user_id" gorm:"index"`
	TargetType    string     `json:"target_type" gorm:"size:20;index"`
	TargetID      uint       `json:"target_id" gorm:"index"`
	Amount        int64      `json:"amount"`
	Message       string     `json:"message" gorm:"size:500"`
	PaymentMethod string     `json:"payment_method" gorm:"size:20"`
	PaymentStatus string     `json:"payment_status" gorm:"size:20;default:pending;index"`
	TradeNo       string     `json:"trade_no" gorm:"size:100;index"`
	PaidAt        *time.Time `json:"paid_at"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}
