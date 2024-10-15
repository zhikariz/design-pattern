package server

import (
	"design-pattern/configs"
	"design-pattern/pkg/route"

	"github.com/labstack/echo/v4"
)

type Server struct {
	*echo.Echo
}

func NewServer(cfg *configs.Config,
	publicRoutes, privateRoutes []route.Route) *Server {
	e := echo.New()
	e.HideBanner = true

	v1 := e.Group("/api/v1")

	if len(publicRoutes) > 0 {
		for _, route := range publicRoutes {
			v1.Add(route.Method, route.Path, route.Handler)
		}
	}

	if len(privateRoutes) > 0 {
		for _, route := range privateRoutes {
			v1.Add(route.Method, route.Path, route.Handler)
		}
	}
	return &Server{e}
}
