package api

import (
	"github.com/gin-gonic/gin"
	"github.com/abdelmounim-dev/go-tshirt/internal/api/handlers"
	"github.com/abdelmounim-dev/go-tshirt/internal/config"
	"github.com/abdelmounim-dev/go-tshirt/internal/db"
	"github.com/abdelmounim-dev/go-tshirt/internal/repository"
	"github.com/abdelmounim-dev/go-tshirt/internal/service"
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
