package builder

import (
	"design-pattern/configs"
	"design-pattern/internal/http/handler"
	"design-pattern/internal/http/router"
	"design-pattern/internal/repository"
	"design-pattern/internal/service"
	"design-pattern/pkg/cache"
	"design-pattern/pkg/route"
	"design-pattern/pkg/token"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func BuildPublicRoutes(cfg *configs.Config, db *gorm.DB, rdb *redis.Client) []route.Route {
	cacheable := cache.NewCacheable(rdb)
	userRepository := repository.NewUserRepository(db)
	tokenUseCase := token.NewTokenUseCase(cfg.JWT.SecretKey)
	userService := service.NewUserService(userRepository, tokenUseCase, cacheable)
	userHandler := handler.NewUserHandler(userService)
	return router.PublicRoutes(userHandler)
}

func BuildPrivateRoutes(cfg *configs.Config, db *gorm.DB, rdb *redis.Client) []route.Route {
	cacheable := cache.NewCacheable(rdb)
	userRepository := repository.NewUserRepository(db)
	tokenUseCase := token.NewTokenUseCase(cfg.JWT.SecretKey)
	userService := service.NewUserService(userRepository, tokenUseCase, cacheable)
	userHandler := handler.NewUserHandler(userService)
	return router.PrivateRoutes(userHandler)
}
