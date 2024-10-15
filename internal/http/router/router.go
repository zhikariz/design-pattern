package router

import (
	"design-pattern/internal/http/handler"
	"design-pattern/pkg/route"
	"net/http"
)

func PublicRoutes(userHandler handler.UserHandler) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: userHandler.FindAll,
		},
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: userHandler.Login,
		},
	}
}

func PrivateRoutes() []route.Route {
	return nil
}
