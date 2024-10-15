package builder

import (
	"design-pattern/internal/http/handler"
	"design-pattern/internal/http/router"
	"design-pattern/internal/repository"
	"design-pattern/internal/service"
	"design-pattern/pkg/route"

	"gorm.io/gorm"
)

func BuildPublicRoutes(db *gorm.DB) []route.Route {
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	return router.PublicRoutes(userHandler)
}

func BuildPrivateRoutes() []route.Route {
	return router.PrivateRoutes()
}
