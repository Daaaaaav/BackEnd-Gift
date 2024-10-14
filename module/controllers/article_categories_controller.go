package controllers

import (
	"module/module/articles"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
)

type ArticleCategoriesController struct {
	DB *gorm.DB
}

func NewArticleCategoriesController(db *gorm.DB) *ArticleCategoriesController {
	return &ArticleCategoriesController{DB: db}
}

func (ac *ArticleCategoriesController) AddCategory(c *gin.Context) {
	var categories articles.ArticleCategories
	if err := c.ShouldBindJSON(&categories); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := bluemonday.UGCPolicy()
	categories.Name = p.Sanitize(categories.Name)
	categories.Slug = p.Sanitize(categories.Slug)
	if err := ac.DB.Create(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category!"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (cc *ArticleCategoriesController) GetCategory(c *gin.Context) {
	id := c.MustGet("id").(int)
	var category articles.ArticleCategories
	if err := cc.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found!"})
		return
	}
	c.JSON(http.StatusOK, category)
}

func (cc *ArticleCategoriesController) GetCategories(ctx *gin.Context) {
	var categories []articles.ArticleCategories
	keyword := ctx.Query("keyword")
	query := cc.DB.Order("created_at desc")
	if keyword != "" {
		query = query.Where("name ILIKE ?", "%"+keyword+"%")
	}
	err := query.Find(&categories).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories!"})
		return
	}
	ctx.JSON(http.StatusOK, categories)
}

func (cc *ArticleCategoriesController) UpdateCategory(c *gin.Context) {
	id := c.MustGet("id").(int)
	var category articles.ArticleCategories
	if err := cc.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found!"})
		return
	}
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := bluemonday.UGCPolicy()
	category.Name = p.Sanitize(category.Name)
	category.Slug = p.Sanitize(category.Slug)
	if err := cc.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category!"})
		return
	}
	c.JSON(http.StatusOK, category)
}

func (cc *ArticleCategoriesController) DeleteCategory(c *gin.Context) {
	id := c.MustGet("id").(int)
	var articleCount int64
	cc.DB.Model(&articles.Articles{}).Where("category_id = ?", id).Count(&articleCount)
	if articleCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete category with associated articles!"})
		return
	}
	if err := cc.DB.Delete(&articles.ArticleCategories{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully!"})
}
package controllers

import (
	"module/module/articles"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
)

type ArticleCategoriesController struct {
	DB *gorm.DB
}

func NewArticleCategoriesController(db *gorm.DB) *ArticleCategoriesController {
	return &ArticleCategoriesController{DB: db}
}

func (ac *ArticleCategoriesController) AddCategory(c *gin.Context) {
	var categories articles.ArticleCategories
	if err := c.ShouldBindJSON(&categories); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := bluemonday.UGCPolicy()
	categories.Name = p.Sanitize(categories.Name)
	categories.Slug = p.Sanitize(categories.Slug)
	if err := ac.DB.Create(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category!"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (cc *ArticleCategoriesController) GetCategory(c *gin.Context) {
	id := c.MustGet("id").(int)
	var category articles.ArticleCategories
	if err := cc.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found!"})
		return
	}
	c.JSON(http.StatusOK, category)
}

func (cc *ArticleCategoriesController) GetCategories(ctx *gin.Context) {
	var categories []articles.ArticleCategories
	keyword := ctx.Query("keyword")
	query := cc.DB.Order("created_at desc")
	if keyword != "" {
		query = query.Where("name ILIKE ?", "%"+keyword+"%")
	}
	err := query.Find(&categories).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories!"})
		return
	}
	ctx.JSON(http.StatusOK, categories)
}

func (cc *ArticleCategoriesController) UpdateCategory(c *gin.Context) {
	id := c.MustGet("id").(int)
	var category articles.ArticleCategories
	if err := cc.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found!"})
		return
	}
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := bluemonday.UGCPolicy()
	category.Name = p.Sanitize(category.Name)
	category.Slug = p.Sanitize(category.Slug)
	if err := cc.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category!"})
		return
	}
	c.JSON(http.StatusOK, category)
}

func (cc *ArticleCategoriesController) DeleteCategory(c *gin.Context) {
	id := c.MustGet("id").(int)
	var articleCount int64
	cc.DB.Model(&articles.Articles{}).Where("category_id = ?", id).Count(&articleCount)
	if articleCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete category with associated articles!"})
		return
	}
	if err := cc.DB.Delete(&articles.ArticleCategories{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully!"})
}
