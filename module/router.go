package module

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"module/module/articles"
	"module/module/controllers"
	"module/module/middlewares"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type APIServer struct {
	listenAddr                  string
	router                      *gin.Engine
	database                    *gorm.DB
	articlesController          *controllers.ArticlesController
	articleCategoriesController *controllers.ArticleCategoriesController
}

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var ctx = context.Background()

func NewAPIServer(listenAddr string, db *gorm.DB) *APIServer {
	router := gin.Default()
	articlesController := controllers.NewArticleController(db)
	categoriesController := controllers.NewArticleCategoriesController(db)
	router.Use(RateLimitMiddleware())

	articles := router.Group("/articles")
	{
		articles.POST("", articlesController.AddArticle)
		articles.GET("", articlesController.GetArticles)
		articles.GET("/:id", middlewares.ValidateIDMiddleware(), articlesController.GetArticle)
		articles.PUT("/:id", middlewares.ValidateIDMiddleware(), articlesController.UpdateArticle)
		articles.DELETE("/:id", middlewares.ValidateIDMiddleware(), articlesController.DeleteArticle)
	}

	categories := router.Group("/categories")
	{
		categories.POST("", categoriesController.AddCategory)
		categories.GET("", categoriesController.GetCategories)
		categories.GET("/:id", middlewares.ValidateIDMiddleware(), categoriesController.GetCategory)
		categories.PUT("/:id", middlewares.ValidateIDMiddleware(), categoriesController.UpdateCategory)
		categories.DELETE("/:id", middlewares.ValidateIDMiddleware(), categoriesController.DeleteCategory)
	}

	return &APIServer{
		listenAddr:                  listenAddr,
		router:                      router,
		database:                    db,
		articlesController:          articlesController,
		articleCategoriesController: categoriesController,
	}
}

var rateLimiter = time.Tick(time.Second / 10)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		<-rateLimiter
		c.Next()
	}
}

func (s *APIServer) GetArticlesWithCache(c *gin.Context) {
	cachedArticles, err := rdb.Get(ctx, "articles").Result()
	if err == redis.Nil {
		var articles []articles.Articles
		err := s.database.Preload("Category").Order("created_at desc").Limit(10).Find(&articles).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch articles! %v", err)})
			return
		}

		jsonData, _ := json.Marshal(articles)
		err = rdb.Set(ctx, "articles", jsonData, time.Hour).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to cache articles! %v", err)})
			return
		}

		c.JSON(http.StatusOK, articles)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Redis error! %v", err)})
		return
	}
	var articles []articles.Articles
	json.Unmarshal([]byte(cachedArticles), &articles)
	c.JSON(http.StatusOK, articles)
}

func (s *APIServer) GetCategoriesWithCache(c *gin.Context) {
	cachedCategories, err := rdb.Get(ctx, "categories").Result()
	if err == redis.Nil {
		var categories []articles.ArticleCategories
		err := s.database.Preload("Articles").Order("created_at desc").Limit(10).Find(&categories).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch categories! %v", err)})
			return
		}
		jsonData, _ := json.Marshal(categories)
		err = rdb.Set(ctx, "categories", jsonData, time.Hour).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to cache categories! %v", err)})
			return
		}
		c.JSON(http.StatusOK, categories)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Redis error! %v", err)})
		return
	}
	var categories []articles.ArticleCategories
	json.Unmarshal([]byte(cachedCategories), &categories)
	c.JSON(http.StatusOK, categories)
}

func (s *APIServer) Run() {
	log.Println("API server running on port", s.listenAddr)
	if err := s.router.Run(s.listenAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
