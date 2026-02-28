package model

type Rescue struct {
	BaseModel
	ReporterID   uint    `json:"reporter_id" gorm:"index"`
	Title        string  `json:"title" gorm:"size:100"`
	Description  string  `json:"description" gorm:"type:text"`
	Species      string  `json:"species" gorm:"size:20"`
	Urgency      string  `json:"urgency" gorm:"size:20;default:medium;index"`
	Status       string  `json:"status" gorm:"size:20;default:reported;index"`
	Longitude    float64 `json:"longitude" gorm:"type:decimal(10,7)"`
	Latitude     float64 `json:"latitude" gorm:"type:decimal(10,7)"`
	Address      string  `json:"address" gorm:"size:500"`
	ContactPhone string  `json:"contact_phone" gorm:"size:20"`

	//关联
	Reporter User           `json:"reporter" gorm:"foreignKey:ReporterID"`
	Images   []RescueImage  `json:"images" gorm:"foreignKey:RescueID"`
	Follows  []RescueFollow `json:"follows" gorm:"foreignKey:RescueID"`
}

type RescueImage struct {
	BaseModel
	RescueID  uint   `json:"rescue_id" gorm:"index"`
	ImageURL  string `json:"image_url" gorm:"size:500"`
	SortOrder int    `json:"sort_order" gorm:"default:0"`
}

type RescueFollow struct {
	BaseModel
	RescueID uint   `json:"rescue_id" gorm:"index"`
	UserID   uint   `json:"user_id" gorm:"index"`
	Content  string `json:"content" gorm:"type:text"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}

type RescueClaim struct {
	BaseModel
	RescueID uint   `json:"rescue_id" gorm:"index"`
	OrgID    uint   `json:"org_id" gorm:"index"`
	Status   string `json:"status" gorm:"size:20;default:claimed"`
	Note     string `json:"note" gorm:"type:text"`

	Rescue       Rescue       `json:"rescue" gorm:"foreignKey:RescueID"`
	Organization Organization `json:"organization" gorm:"foreignKey:OrgID"`
}
