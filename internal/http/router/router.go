package router

import (
	"design-pattern/internal/http/handler"
	"design-pattern/pkg/route"
	"net/http"
)

func PublicRoutes(userHandler handler.UserHandler) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: userHandler.Login,
		},
		{
			Method:  http.MethodGet,
			Path:    "/generate-password/:password",
			Handler: userHandler.GeneratePassword,
		},
	}
}

func PrivateRoutes(userHandler handler.UserHandler) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: userHandler.FindAll,
			Roles:   []string{"admin", "editor"},
		},
	}
}
