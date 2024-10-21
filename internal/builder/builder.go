package builder

import (
	"design-pattern/configs"
	"design-pattern/internal/http/handler"
	"design-pattern/internal/http/router"
	"design-pattern/internal/repository"
	"design-pattern/internal/service"
	"design-pattern/pkg/route"
	"design-pattern/pkg/token"

	"gorm.io/gorm"
)

func BuildPublicRoutes(cfg *configs.Config, db *gorm.DB) []route.Route {
	userRepository := repository.NewUserRepository(db)
	tokenUseCase := token.NewTokenUseCase(cfg.JWT.SecretKey)
	userService := service.NewUserService(userRepository, tokenUseCase)
	userHandler := handler.NewUserHandler(userService)
	return router.PublicRoutes(userHandler)
}

func BuildPrivateRoutes(cfg *configs.Config, db *gorm.DB) []route.Route {
	userRepository := repository.NewUserRepository(db)
	tokenUseCase := token.NewTokenUseCase(cfg.JWT.SecretKey)
	userService := service.NewUserService(userRepository, tokenUseCase)
	userHandler := handler.NewUserHandler(userService)
	return router.PrivateRoutes(userHandler)
}
