package server

import (
	"design-pattern/configs"

	"github.com/labstack/echo/v4"
)

type Server struct {
	*echo.Echo
}

func NewServer(cfg *configs.Config) *Server {
	e := echo.New()
	e.HideBanner = true
	return &Server{e}
}
