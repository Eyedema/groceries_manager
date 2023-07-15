package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"ubaldo/api_server/internal/handlers"
	"ubaldo/api_server/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/patrickmn/go-cache"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var itemCache *cache.Cache

func NewServer() http.Handler {
	var err error
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=5432 sslmode=disable TimeZone=UTC",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	err = db.AutoMigrate(&models.Item{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}
	itemCache = cache.New(cache.DefaultExpiration, cache.DefaultExpiration)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/is-alive"))

	r.Get("/items", handlers.GetAllItems)
	r.Get("/item/{id}", handlers.GetItemByID)
	r.Post("/item", handlers.SaveItem)
	r.Delete("/item/{id}", handlers.DeleteItem)

	return r
}
