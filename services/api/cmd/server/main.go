package main

import (
	"github.com/vaxxnsh/metaverse/api/internal/config"
	"github.com/vaxxnsh/metaverse/api/internal/repository"
	"github.com/vaxxnsh/metaverse/api/internal/service"
)

func main() {
	cfg := config.Load()

	db, queries, _ := db.NewDB(cfg.DBURL)

	userRepo := repository.NewUserRepository(queries)
	userService := service.NewUserService()
	userHandler := handlers.NewUserHandler(userService)

	router := router.SetupRouter(userHandler)

	router.Run(cfg.Port)
}
