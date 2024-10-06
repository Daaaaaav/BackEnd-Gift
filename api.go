package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/your_module_name/controllers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type APIServer struct {
	listenAddr string
	router     *gin.Engine
	database   *gorm.DB
}

func NewAPIServer(listenAddr string, db *gorm.DB) *APIServer {
	router := gin.Default()
	router.Use(RateLimitMiddleware())

	return &APIServer{
		listenAddr: listenAddr,
		router:     router,
		database:   db,
	}
}

func (s *APIServer) Run() {
	articleController := controllers.NewArticleController(s.database)

	// Article Routes
	s.router.POST("/articles", articleController.AddArticle)
	s.router.GET("/articles", articleController.GetArticles)
	s.router.GET("/articles/:id", articleController.GetArticle)
	s.router.PUT("/articles/:id", articleController.UpdateArticle)
	s.router.DELETE("/articles/:id", articleController.DeleteArticle)

	log.Println("API server running on port", s.listenAddr)
	http.ListenAndServe(s.listenAddr, s.router)
}

var rateLimiter = time.Tick(time.Second / 10)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		<-rateLimiter
		c.Next()
	}
}

func main() {
	// Connect to the database (PostgreSQL example)
	dsn := "host=localhost user=postgres dbname=godb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create API server
	server := NewAPIServer(":8080", db)
	server.Run()
}
