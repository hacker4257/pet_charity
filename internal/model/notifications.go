package model

type Notification struct {
	BaseModel
	UserID    uint   `json:"user_id" gorm:"index;not null"`
	Type      string `json:"type" gorm:"size:30;not null"`
	Title     string `json:"title" gorm:"size:100;not null"`
	Content   string `json:"content" gorm:"size:500"`
	RelatedID uint   `json:"related_id" gorm:"default:0"`
	IsRead    bool   `json:"is_read" gorm:"default:false;index"`
	User      User   `json:"-" gorm:"foreignKey:UserID"`
}
