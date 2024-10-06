package articles

type Articles struct {
	ID         uint              `gorm:"primaryKey" json:"id"`
	CategoryID uint              `json:"category_id"`
	Title      string            `gorm:"size:100;not null" json:"title"`
	Content    string            `gorm:"type:text;not null" json:"content"`
	Thumbnail  string            `gorm:"size:100;not null" json:"thumbnail"`
	Status     bool              `gorm:"default:true" json:"status"`
	Slug       string            `gorm:"size:255;not null;unique" json:"slug"`
	Category   ArticleCategories `gorm:"foreignKey:CategoryID" json:"category"`
}
