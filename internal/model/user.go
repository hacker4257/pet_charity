package model

type User struct {
	BaseModel
	Username      string `json:"username" gorm:"size:50;uniqueIndex"`
	Email         string `json:"email" gorm:"size:100;uniqueIndex"`
	Phone         string `json:"phone" gorm:"size:20;index"`
	Password      string `json:"-" gorm:"size:255"`
	Avatar        string `json:"avatar" gorm:"size:255"`
	Nickname      string `json:"nickname" gorm:"size:50"`
	Role          string `json:"role" gorm:"size:20;default:user;index"`
	Status        string `json:"status" gorm:"size:20;default:active"`
	Language      string `json:"language" gorm:"size:10;default:zh-CN"`
	ActivityScore int    `json:"activity_score" gorm:"default:0;index"`
}

type UserActivityLogs struct {
	BaseModel
	UserID uint   `json:"user_id" gorm:"index"`
	Action string `json:"action" gorm:"size:30;index"` // login, donation, adoption, rescue, chat
	Points int    `json:"points"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}
