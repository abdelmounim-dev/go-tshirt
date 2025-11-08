package api

import (
	"github.com/gin-gonic/gin"
	"go-tshirt/internal/config"
	"go-tshirt/internal/db"
	"go-tshirt/internal/repository"
	"go-tshirt/internal/service"
	"go-tshirt/internal/api/handlers"
)

func NewRouter(cfg config.Config) *gin.Engine {
	r := gin.Default()

	database, err := db.NewSQLite(cfg.DBPath)
	if err != nil {
		panic(err)
	}

	repo := repository.NewProductRepo(database)
	_ = repo.AutoMigrate()

	svc := service.NewProductService(repo)
	handler := handlers.NewProductHandler(svc)

	api := r.Group("/api/products")
	handler.Register(api)

	return r
}
