package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vaxxnsh/metaverse/api/internal/config"
	"github.com/vaxxnsh/shared/db"
	"github.com/vaxxnsh/metaverse/api/internal/handler"
	"github.com/vaxxnsh/metaverse/api/internal/repository"
	"github.com/vaxxnsh/metaverse/api/internal/router"
	"github.com/vaxxnsh/metaverse/api/internal/service"
)

func NewDB(dbURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

func main() {
	cfg := config.Load()

	pool, err := NewDB(cfg.DBURL)
	if err != nil {
		log.Fatal("error while connecting with database")
	}
	queries := db.New(pool)
	userRepo := repository.NewUserRepository(queries)
	userService := service.NewUserService(userRepo)
	adminRepo := repository.NewAdminRepository(queries)
	adminService := service.NewAdminService(adminRepo)
	metadataRepo := repository.NewMetadataRepository(queries)
	metadataService := service.NewMetadataService(metadataRepo, userService)
	spaceRepo := repository.NewSpaceRepository(pool, queries)
	spaceService := service.NewSpaceService(spaceRepo)
	arenaRepo := repository.NewArenaRepository(queries)
	arenaService := service.NewArenaService(arenaRepo)
	mapCreatorRepo := repository.NewMapCreatorRepository(pool, queries)
	mapCreatorService := service.NewMapCreatorService(mapCreatorRepo)

	authService := service.NewAuthService(userService, adminService)

	appHandler := router.AppHandlers{
		AuthHandler:     *handler.NewAuthHandler(authService),
		UserHandler:     *handler.NewUserHandler(userService),
		AdminHandler:    *handler.NewAdminHandler(adminService),
		MetadataHandler: *handler.NewMetadataHandler(metadataService),
		SpaceHandler:    *handler.NewSpaceHandler(spaceService),
		ArenaHandler:       *handler.NewArenaHandler(arenaService),
		MapCreatorHandler:  *handler.NewMapCreatorHandler(mapCreatorService),
	}

	router := router.SetupRouter(appHandler)
	router.Run(fmt.Sprintf(":%s", cfg.Port))
}
