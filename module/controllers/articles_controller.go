package controllers

import (
	"module/module/articles"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
)

type ArticlesController struct {
	DB *gorm.DB
}

func NewArticleController(db *gorm.DB) *ArticlesController {
	return &ArticlesController{DB: db}
}

func (ac *ArticlesController) AddArticle(c *gin.Context) {
	var article articles.Articles
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := bluemonday.UGCPolicy()
	article.Title = p.Sanitize(article.Title)
	article.Content = p.Sanitize(article.Content)
	article.Thumbnail = p.Sanitize(article.Thumbnail)
	article.Slug = p.Sanitize(article.Slug)
	if err := ac.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article!"})
		return
	}
	c.JSON(http.StatusOK, article)
}

func (ac *ArticlesController) GetArticle(c *gin.Context) {
	id := c.Param("id")
	var article articles.Articles
	if err := ac.DB.Preload("Category").First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found!"})
		return
	}
	c.JSON(http.StatusOK, article)
}

func (ac *ArticlesController) GetArticles(ctx *gin.Context) {
	var articles []articles.Articles
	keyword := ctx.Query("keyword")
	categoryID := ctx.Query("category_id")
	query := ac.DB.Preload("Category").Order("created_at desc")
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		query = query.Where("title ILIKE ?", "%"+keyword+"%")
	}
	err := query.Find(&articles).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles!"})
		return
	}
	ctx.JSON(http.StatusOK, articles)
}

func (ac *ArticlesController) UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	var article articles.Articles
	if err := ac.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found!"})
		return
	}
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := bluemonday.UGCPolicy()
	article.Title = p.Sanitize(article.Title)
	article.Content = p.Sanitize(article.Content)
	if err := ac.DB.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article!"})
		return
	}
	c.JSON(http.StatusOK, article)
}

func (ac *ArticlesController) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID, try again!"})
		return
	}
	if err := ac.DB.Delete(&articles.Articles{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Article successful to be deleted!"})
}
