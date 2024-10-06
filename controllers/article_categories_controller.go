package controllers

import (
	"e/module/articles"
	"e/module/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateAnArticleCategory(c *gin.Context) {
	var categories articles.ArticleCategories
	if err := c.ShouldBindJSON(&categories); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := data.DB.Create(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}
