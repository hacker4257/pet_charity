package model

type Pet struct {
	BaseModel
	OrgID        uint   `json:"org_id" gorm:"index"`
	Name         string `json:"name" gorm:"size:50"`
	Species      string `json:"species" gorm:"size:20;index"`
	Breed        string `json:"breed" gorm:"size:50"`
	Age          int    `json:"age"`
	Gender       string `json:"gender" gorm:"size:10"`
	HealthStatus string `json:"health_status" gorm:"size:50"`
	Description  string `json:"description" gorm:"type:text"`
	CoverImage   string `json:"cover_image" gorm:"size:500"`
	Status       string `json:"status" gorm:"size:20;default:available;index"`

	Organization Organization `json:"organization" gorm:"foreignKey:OrgID"`
	Images       []PetImage   `json:"images" gorm:"foreignKey:PetID"`
}

type PetImage struct {
	BaseModel
	PetID     uint   `json:"pet_id" gorm:"index"`
	ImageURL  string `json:"image_url" gorm:"size:500"`
	SortOrder int    `json:"sort_order" gorm:"default:0"`
}
