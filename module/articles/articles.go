package articles

import (
	"time"
)

type Articles struct {
	ID         uint              `json:"id" gorm:"primaryKey"`
	Title      string            `json:"title"`
	Content    string            `json:"content"`
	Thumbnail  string            `json:"thumbnail"`
	Status     bool              `json:"status"`
	Slug       string            `json:"slug"`
	CategoryID uint              `json:"category_id"`
	Category   ArticleCategories `json:"category" gorm:"constraint:OnDelete:CASCADE;"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}
