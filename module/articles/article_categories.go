package articles

type ArticleCategories struct {
	ID       uint       `gorm:"not null;primaryKey" json:"id"`
	Name     string     `gorm:"not null;unique" json:"name"`
	Slug     string     `gorm:"not null;unique" json:"slug"`
	Articles []Articles `gorm:"foreignKey:CategoryID"`
}
