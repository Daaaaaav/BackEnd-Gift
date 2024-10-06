package main

import (
	"context"
	"e/module/controllers"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Articles struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Thumbnail string    `json:"thumbnail"`
	Status    bool      `json:"status"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ArticlesController struct {
	DB *gorm.DB
}

func NewArticleController(db *gorm.DB) *ArticlesController {
	return &ArticlesController{DB: db}
}

type APIServer struct {
	listenAddr         string
	router             *gin.Engine
	database           *gorm.DB
	articlesController *controllers.ArticlesController
}

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var ctx = context.Background()

func (ac *ArticlesController) AddArticle(c *gin.Context) {
	var article Articles
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	c.JSON(http.StatusOK, article)
}

func NewAPIServer(listenAddr string, db *gorm.DB) *APIServer {
	router := gin.Default()
	articlesController := controllers.NewArticleController(db)
	router.Use(RateLimitMiddleware())
	return &APIServer{
		listenAddr:         listenAddr,
		router:             router,
		database:           db,
		articlesController: articlesController,
	}
}

var rateLimiter = time.Tick(time.Second / 10)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		<-rateLimiter
		c.Next()
	}
}

func (s *APIServer) GetArticles(c *gin.Context) {
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var articles []Articles
	err = s.database.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&articles).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch articles: %v", err)})
		return
	}

	c.JSON(http.StatusOK, articles)
}

func (s *APIServer) GetArticlesWithCache(c *gin.Context) {
	cachedArticles, err := rdb.Get(ctx, "articles").Result()
	if err == redis.Nil {
		var articles []Articles
		err := s.database.Preload("Category").Order("created_at desc").Limit(10).Find(&articles).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch articles: %v", err)})
			return
		}

		jsonData, _ := json.Marshal(articles)
		err = rdb.Set(ctx, "articles", jsonData, time.Hour).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to cache articles: %v", err)})
			return
		}

		c.JSON(http.StatusOK, articles)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("redis error: %v", err)})
		return
	}
	var articles []Articles
	json.Unmarshal([]byte(cachedArticles), &articles)
	c.JSON(http.StatusOK, articles)
}

func (s *APIServer) Run() {
	s.router.POST("/articles", s.articlesController.AddArticle)
	s.router.GET("/articles", s.articlesController.GetArticles)
	s.router.GET("/articles/:id", s.articlesController.GetArticle)
	s.router.PUT("/articles/:id", s.articlesController.UpdateArticle)
	s.router.DELETE("/articles/:id", s.articlesController.DeleteArticle)
	log.Println("API server running on port", s.listenAddr)
	http.ListenAndServe(s.listenAddr, s.router)
}

func main() {
	dsn := "host=localhost user=postgres  password=Davina241105 dbname=dbarticle port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	server := NewAPIServer(":8080", db)
	server.Run()
}
