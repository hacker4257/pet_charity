package model


type PetDiary struct {
	BaseModel
	UserID uint `json:"user_id" gorm:"index"`
	PetID uint `json:"pet_id" gorm:"index"`
	Content string `json:"content" gorm:"type:text"`

	//关联表
	User User `json:"user" gorm:"foreignKey:UserID"`
	Pet Pet `json:"pet" gorm:"foreignKey:PetID"`
	Images []DiaryImage `json:"images" gorm:"foreignKey:DiaryID"`
}

type DiaryImage struct {
	BaseModel

	DiaryID uint `json:"diary_id" gorm:"index"`
	ImageURL string `json:"image_url" gorm:"size:500"`
	SortOrder int `json:"sort_order" gorm:"default:0"`
}

type DiaryLike struct {
	BaseModel

	DiaryID uint `json:"diary_id" gorm:"uniqueIndex:idx_diary_user"`
	UserID uint `json:"user_id" gorm:"uniqueIndex:idx_diary_user"`
}