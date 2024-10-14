package controllers

import (
	"module/module/articles"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

	if err := ac.DB.Create(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article!"})
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

func (ac *ArticleCategoriesController) GetCategories(c *gin.Context) {
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	var categories []articles.ArticleCategories
	if err := ac.DB.Preload("Category").Order("created_at desc").Offset(offset).Limit(pageSize).Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch!"})
		return
	}
	c.JSON(http.StatusOK, categories)
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
